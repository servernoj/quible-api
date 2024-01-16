package BasketAPI

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/quible-io/quible-api/lib/models"
)

func getTeamEnhancer(ctx context.Context) (func(MatchScheduleTeam) TeamInfo, error) {
	// teams info
	teamsInfo, err := models.TeamInfos().AllG(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve combined team info (for all teams): %w", err)
	}
	teamInfoById := make(map[int]*models.TeamInfo, len(teamsInfo))
	for _, teamInfo := range teamsInfo {
		teamInfoById[teamInfo.ID] = teamInfo
	}
	return func(team MatchScheduleTeam) TeamInfo {
		teamInfo := *teamInfoById[int(team.ID)]
		return TeamInfo{
			ID:             teamInfo.ID,
			Name:           teamInfo.Name,
			Slug:           teamInfo.Slug,
			ShortName:      teamInfo.ShortName,
			Abbr:           teamInfo.Abbr,
			ArenaName:      teamInfo.ArenaName,
			ArenaSize:      teamInfo.ArenaSize,
			Color:          teamInfo.Color,
			SecondaryColor: teamInfo.SecondaryColor,
			Logo:           teamInfo.Logo.Ptr(),
		}
	}, nil
}

func GetGames(ctx context.Context, query GetGamesDTO) ([]Game, error) {
	loc, _ := time.LoadLocation("America/New_York")
	if query.LocalTimeZoneShift != nil {
		loc = time.FixedZone("User timezone", *query.LocalTimeZoneShift*int(time.Hour/time.Second))
	}
	dateParsed, _ := time.Parse(time.DateOnly, query.Date)
	dateParsedInLocation, _ := time.ParseInLocation(time.DateOnly, query.Date, loc)
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
	// -- inject team details into games
	teamEnhancer, err := getTeamEnhancer(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize team entity enhancer: %w", err)
	}
	games := make([]Game, 0, len(body.Events))
	tsFrom := dateParsedInLocation.Unix()
	tsTo := dateParsedInLocation.Add(24 * time.Hour).Unix()
	// -- pre-populate `games` slice
	for _, ev := range body.Events {
		if ev.Tournament.Name != "NBA" || ev.StartTimestamp < tsFrom || ev.StartTimestamp > tsTo {
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
		game := Game{
			ID:         ev.ID,
			GameStatus: GameStatus,
			AwayScore:  ev.AwayScore.Current,
			HomeScore:  ev.HomeScore.Current,
			Date:       time.Unix(ev.StartTimestamp, 0).Format(time.RFC3339),
			HomeTeam:   teamEnhancer(ev.HomeTeam),
			AwayTeam:   teamEnhancer(ev.AwayTeam),
		}
		games = append(games, game)
	}
	return games, nil
}
