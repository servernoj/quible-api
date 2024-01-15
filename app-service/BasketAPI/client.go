package BasketAPI

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/quible-io/quible-api/lib/models"
)

func getTeamEnhancer(ctx context.Context, teamDetailsMap map[uint]*TeamDetails) (func(MatchScheduleTeam) TeamInfo, error) {
	// logos
	images, err := models.Images(models.ImageWhere.ParentID.IsNull()).AllG(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve logos for all teams: %w", err)
	}
	logoByShortName := make(map[string]*string, len(images))
	for _, image := range images {
		imageUrl := image.ImageURL
		logoByShortName[image.DisplayName] = &imageUrl
	}
	return func(team MatchScheduleTeam) TeamInfo {
		teamInfo := TeamInfo{
			ID:        team.ID,
			Name:      team.Name,
			Slug:      team.Slug,
			ShortName: team.ShortName,
			Abbr:      team.NameCode,
			Logo:      logoByShortName[team.ShortName],
		}
		if teamDetails, ok := teamDetailsMap[team.ID]; ok {
			teamInfo.Color = teamDetails.TeamColors.Primary
			teamInfo.SecondaryColor = teamDetails.TeamColors.Secondary
			teamInfo.ArenaName = teamDetails.Venue.Stadium.Name
			teamInfo.ArenaSize = teamDetails.Venue.Stadium.Capacity
		} else {
			fmt.Printf("no details for team %d\n", team.ID)
		}
		return teamInfo
	}, nil
}

func GetGames(ctx context.Context, date string) ([]Game, error) {
	dateParsed, err := time.Parse(time.DateOnly, date)
	if err != nil || dateParsed.Format(time.DateOnly) != date {
		return nil, fmt.Errorf("unable to [correctly] parse date %q as YYYY-MM-DD: %w", date, err)
	}
	host := "basketapi1.p.rapidapi.com"
	url := fmt.Sprintf("https://%s/api/basketball/matches/%s", host, dateParsed.Format("2/1/2006"))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-RapidAPI-Key", os.Getenv("ENV_RAPIDAPI_KEY"))
	req.Header.Add("X-RapidAPI-Host", host)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable execute request to %q: %w", url, err)
	}
	var body MatchScheduleData
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("unable to decode response: %w", err)
	}
	res.Body.Close()
	games := make([]Game, 0, len(body.Events))
	teamDetailsById := make(map[uint]*TeamDetails)
	// -- pre-populate `games` slice
	for _, ev := range body.Events {
		if ev.Tournament.Name != "NBA" {
			continue
		}
		ev := ev
		GameStatus := ev.Status.Description
		if ev.Status.Type == Inprogress && ev.Time.Played != nil {
			totalSeconds := *ev.Time.Played - *ev.Time.PeriodLength**ev.Time.TotalPeriodCount
			if totalSeconds <= 0 {
				totalSeconds = *ev.Time.Played % *ev.Time.PeriodLength
			}
			minutes := totalSeconds / 60
			seconds := totalSeconds % 60
			GameStatus = fmt.Sprintf("%s (%d:%02d)", GameStatus, minutes, seconds)
		}
		teamDetailsById[ev.HomeTeam.ID] = nil
		teamDetailsById[ev.AwayTeam.ID] = nil
		game := Game{
			ID:         ev.ID,
			GameStatus: GameStatus,
			AwayScore:  ev.AwayScore.Current,
			HomeScore:  ev.HomeScore.Current,
			Date:       time.Unix(ev.StartTimestamp, 0).Format(time.RFC3339),
			Event:      &ev,
		}
		games = append(games, game)
	}
	// -- fetch team details on referenced teams
	// producer
	teamIds := make(chan uint)
	teamDetails := make(chan *TeamDetails, len(teamDetailsById))
	concurrency := 4
	RPS := 10
	ticker := time.NewTicker(time.Second / time.Duration(RPS*concurrency))
	throttled := make(chan uint)
	var wg sync.WaitGroup
	wg.Add(concurrency)
	successCounter := atomic.Int64{}
	go func() {
		for id := range teamDetailsById {
			teamIds <- id
		}
	}()
	// throttler
	go func() {
		for id := range teamIds {
			<-ticker.C
			throttled <- id
		}
		ticker.Stop()
		close(throttled)
	}()
	// concurrent executor
	for i := 0; i < concurrency; i++ {
		go func() {
			for id := range throttled {
				url := fmt.Sprintf("https://%s/api/basketball/team/%d", host, id)
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					log.Printf("unable to create team details request: %s", err)
				}
				req.Header.Add("X-RapidAPI-Key", os.Getenv("ENV_RAPIDAPI_KEY"))
				req.Header.Add("X-RapidAPI-Host", host)
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Printf("unable to execute team details request: %s", err)
					teamIds <- id
					continue
				}
				if res.StatusCode != http.StatusOK {
					log.Printf("[%d] http request failed: %s", id, res.Status)
					teamIds <- id
					continue
				}
				var body TeamDetailsData
				if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
					log.Printf("unable to decode team details response: %s", err)
				}
				teamDetails <- &body.Team
				res.Body.Close()
				// -- register `id` as successful
				successCounter.Add(1)
				if int(successCounter.Load()) >= len(teamDetailsById) {
					close(teamIds)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(teamDetails)
	// load responses into teamDetailsById
	for details := range teamDetails {
		teamDetailsById[details.ID] = details
	}
	// inject team details into games
	teamEnhancer, err := getTeamEnhancer(ctx, teamDetailsById)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize team entity enhancer: %w", err)
	}
	for idx := range games {
		game := &games[idx]
		game.HomeTeam = teamEnhancer(game.Event.HomeTeam)
		game.AwayTeam = teamEnhancer(game.Event.AwayTeam)
	}
	return games, nil
}
