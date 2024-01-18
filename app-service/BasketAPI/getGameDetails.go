package BasketAPI

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/quible-io/quible-api/lib/misc"
)

func getTeamStats(query url.Values) (*GameTeamsStats, error) {
	url := fmt.Sprintf(
		"https://%s/api/basketball/match/%s/statistics",
		Host,
		query.Get("gameId"),
	)
	response, err := misc.GetOne[MatchStatisticsEntries]{
		Client: *http.DefaultClient,
		URL:    url,
		UpdateRequest: func(req *http.Request) {
			req.Header.Set("X-RapidAPI-Key", os.Getenv("ENV_RAPIDAPI_KEY"))
			req.Header.Set("X-RapidAPI-Host", Host)
		},
		ExpectedStatus: http.StatusOK,
	}.Do()
	if err != nil {
		return nil, fmt.Errorf("unable to get response of MatchStatistics: %w", err)
	}
	var statGroups []StatisticsEntryGroup
	for idx := range response.Statistics {
		if response.Statistics[idx].Period == "ALL" {
			statGroups = response.Statistics[idx].Groups
			break
		}
	}
	if len(statGroups) == 0 {
		return nil, fmt.Errorf("statistics groups not found")
	}
	var statItems []GroupItem
	for idx := range statGroups {
		if statGroups[idx].GroupName == GroupName_Other {
			statItems = statGroups[idx].StatisticsItems
			break
		}
	}
	if len(statItems) == 0 {
		return nil, fmt.Errorf("statistics group items not found")
	}
	var result GameTeamsStats
	for _, item := range statItems {
		switch item.Name {
		case Other_Assists:
			{
				result.HomeTeam.Assists = item.HomeValue
				result.AwayTeam.Assists = item.AwayValue
			}

		case Other_Blocks:
			{
				result.HomeTeam.Blocks = item.HomeValue
				result.AwayTeam.Blocks = item.AwayValue
			}
		case Other_Fouls:
			{
				result.HomeTeam.Fouls = item.HomeValue
				result.AwayTeam.Fouls = item.AwayValue
			}
		case Other_Rebounds:
			{
				result.HomeTeam.Rebounds = item.HomeValue
				result.AwayTeam.Rebounds = item.AwayValue
			}
		case Other_Steals:
			{
				result.HomeTeam.Steals = item.HomeValue
				result.AwayTeam.Steals = item.AwayValue
			}
		case Other_Turnovers:
			{
				result.HomeTeam.Turnovers = item.HomeValue
				result.AwayTeam.Turnovers = item.AwayValue
			}
		}
	}

	return &result, nil
}

func getMatchDetails(query url.Values) (*MatchDetails, error) {
	url := fmt.Sprintf(
		"https://%s/api/basketball/match/%s",
		Host,
		query.Get("gameId"),
	)
	response, err := misc.GetOne[MatchData]{
		Client: *http.DefaultClient,
		URL:    url,
		UpdateRequest: func(req *http.Request) {
			req.Header.Set("X-RapidAPI-Key", os.Getenv("ENV_RAPIDAPI_KEY"))
			req.Header.Set("X-RapidAPI-Host", Host)
		},
		ExpectedStatus: http.StatusOK,
	}.Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve match event: %w", err)
	}
	// enhance response
	ev := response.Event
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
	return &MatchDetails{
		ID:         ev.ID,
		GameStatus: GameStatus,
		Date:       time.Unix(ev.StartTimestamp, 0).Format(time.RFC3339),
		AwayScore:  ev.AwayScore.Current,
		HomeScore:  ev.HomeScore.Current,
		Event:      &ev,
	}, nil
}

func GetGameDetails(ctx context.Context, query url.Values) (*GameDetails, error) {
	teamEnhancer, err := getTeamEnhancer(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize team entity enhancer: %w", err)
	}
	matchDetails, err := getMatchDetails(query)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve match details: %w", err)
	}
	teamsStats, err := getTeamStats(query)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve teams stats injector: %w", err)
	}

	return &GameDetails{
		MatchDetails: *matchDetails,
		HomeTeam: TeamInfoExtended{
			TeamInfo: teamEnhancer(matchDetails.Event.HomeTeam),
			Stats:    teamsStats.HomeTeam,
		},
		AwayTeam: TeamInfoExtended{
			TeamInfo: teamEnhancer(matchDetails.Event.AwayTeam),
			Stats:    teamsStats.AwayTeam,
		},
	}, nil
}
