package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/quible-io/quible-api/app-service-ng/services/BasketAPI"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/misc"
)

type GetGameDetailsInput struct {
	AuthorizationHeaderResolver
	GameId uint `query:"gameId"`
}

type GetGameDetailsOutput struct {
	Body GameDetails
}

func ApplyMapper[F any, T any](s []F, m func(F) T) []T {
	result := make([]T, len(s))
	for idx := range s {
		result[idx] = m(s[idx])
	}
	return result
}

func getTeamStats(gameId uint) (*GameTeamsStats, error) {
	url := fmt.Sprintf(
		"https://%s/api/basketball/match/%d/statistics",
		BasketAPI.Host,
		gameId,
	)
	response, err := misc.GetOne[MStat_Data]{
		Client: *http.DefaultClient,
		URL:    url,
		UpdateRequest: func(req *http.Request) {
			req.Header.Set("X-RapidAPI-Key", os.Getenv("ENV_RAPIDAPI_KEY"))
			req.Header.Set("X-RapidAPI-Host", BasketAPI.Host)
		},
		ExpectedStatus: http.StatusOK,
	}.Do()
	if err != nil {
		return nil, err
	}
	var statGroups []MStat_Group
	for idx := range response.Statistics {
		if response.Statistics[idx].Period == "ALL" {
			statGroups = response.Statistics[idx].Groups
			break
		}
	}
	if len(statGroups) == 0 {
		return nil, errors.New("statistics groups not found")
	}
	var statItems []MStat_GroupItem
	var result GameTeamsStats
	for idx := range statGroups {
		isGroupOthers := statGroups[idx].GroupName == MStat_GroupName_Other
		isGroupScoring := statGroups[idx].GroupName == MStat_GroupName_Scoring
		if isGroupOthers || isGroupScoring {
			statItems = statGroups[idx].StatisticsItems
			if len(statItems) == 0 {
				return nil, errors.New("statistics group items not found")
			}
			for _, item := range statItems {
				switch item.Name {
				case MStat_GroupItemName_ScoringFieldGoals:
					{
						if item.HomeTotal != nil {
							result.HomeTeam.FieldGoalAttempts = *item.HomeTotal
						}
						if item.AwayTotal != nil {
							result.AwayTeam.FieldGoalAttempts = *item.AwayTotal
						}
						result.HomeTeam.FieldGoalsMade = item.HomeValue
						result.AwayTeam.FieldGoalsMade = item.AwayValue
					}
				case MStat_GroupItemName_ScoringFreeThrows:
					{
						if item.HomeTotal != nil {
							result.HomeTeam.FreeThrowAttempts = *item.HomeTotal
						}
						if item.AwayTotal != nil {
							result.AwayTeam.FreeThrowAttempts = *item.AwayTotal
						}
						result.HomeTeam.FreeThrowsMade = item.HomeValue
						result.AwayTeam.FreeThrowsMade = item.AwayValue
					}
				case MStat_GroupItemName_ScoringThreePoints:
					{
						if item.HomeTotal != nil {
							result.HomeTeam.ThreePointAttempts = *item.HomeTotal
						}
						if item.AwayTotal != nil {
							result.AwayTeam.ThreePointAttempts = *item.AwayTotal
						}
						result.HomeTeam.ThreePointsMade = item.HomeValue
						result.AwayTeam.ThreePointsMade = item.AwayValue
					}
				case MStat_GroupItemName_OtherAssists:
					{
						result.HomeTeam.Assists = item.HomeValue
						result.AwayTeam.Assists = item.AwayValue
					}
				case MStat_GroupItemName_OtherBlocks:
					{
						result.HomeTeam.Blocks = item.HomeValue
						result.AwayTeam.Blocks = item.AwayValue
					}
				case MStat_GroupItemName_OtherFouls:
					{
						result.HomeTeam.Fouls = item.HomeValue
						result.AwayTeam.Fouls = item.AwayValue
					}
				case MStat_GroupItemName_OtherRebounds:
					{
						result.HomeTeam.Rebounds = item.HomeValue
						result.AwayTeam.Rebounds = item.AwayValue
					}
				case MStat_GroupItemName_OtherSteals:
					{
						result.HomeTeam.Steals = item.HomeValue
						result.AwayTeam.Steals = item.AwayValue
					}
				case MStat_GroupItemName_OtherTurnovers:
					{
						result.HomeTeam.Turnovers = item.HomeValue
						result.AwayTeam.Turnovers = item.AwayValue
					}
				}
			}
		}
	}
	return &result, nil
}

func getPlayersStats(gameId uint) (*GamePlayers, error) {
	url := fmt.Sprintf(
		"https://%s/api/basketball/match/%d/lineups",
		BasketAPI.Host,
		gameId,
	)
	response, err := misc.GetOne[ML_Data]{
		Client: *http.DefaultClient,
		URL:    url,
		UpdateRequest: func(req *http.Request) {
			req.Header.Set("X-RapidAPI-Key", os.Getenv("ENV_RAPIDAPI_KEY"))
			req.Header.Set("X-RapidAPI-Host", BasketAPI.Host)
		},
		ExpectedStatus: http.StatusOK,
	}.Do()
	if err != nil {
		return nil, err
	}
	mapper := func(playerElement ML_PlayerElement) PlayerEntity {
		return PlayerEntity{
			ID:   playerElement.Player.ID,
			Name: playerElement.Player.Name,
			Stats: PlayerStats{
				MinutesPlayed:      float64(playerElement.Statistics.SecondsPlayed) / 60.0,
				SecondsPlayed:      playerElement.Statistics.SecondsPlayed,
				FieldGoalsMade:     playerElement.Statistics.FieldGoalsMade,
				FieldGoalAttempts:  playerElement.Statistics.FieldGoalAttempts,
				ThreePointsMade:    playerElement.Statistics.ThreePointsMade,
				ThreePointAttempts: playerElement.Statistics.FreeThrowAttempts,
				FreeThrowsMade:     playerElement.Statistics.FreeThrowsMade,
				FreeThrowAttempts:  playerElement.Statistics.FreeThrowAttempts,
				OffensiveRebounds:  playerElement.Statistics.OffensiveRebounds,
				DefensiveRebounds:  playerElement.Statistics.DefensiveRebounds,
				Rebounds:           playerElement.Statistics.Rebounds,
				Assists:            playerElement.Statistics.Assists,
				Steals:             playerElement.Statistics.Steals,
				Blocks:             playerElement.Statistics.Blocks,
				Turnovers:          playerElement.Statistics.Turnovers,
				PersonalFouls:      playerElement.Statistics.PersonalFouls,
				Points:             playerElement.Statistics.Points,
			},
		}
	}
	return &GamePlayers{
		HomeTeam: ApplyMapper(response.Home.Players, mapper),
		AwayTeam: ApplyMapper(response.Away.Players, mapper),
	}, nil
}

func getMatchDetails(gameId uint) (*MatchDetails, error) {
	url := fmt.Sprintf(
		"https://%s/api/basketball/match/%d",
		BasketAPI.Host,
		gameId,
	)
	response, err := misc.GetOne[MD_Data]{
		Client: *http.DefaultClient,
		URL:    url,
		UpdateRequest: func(req *http.Request) {
			req.Header.Set("X-RapidAPI-Key", os.Getenv("ENV_RAPIDAPI_KEY"))
			req.Header.Set("X-RapidAPI-Host", BasketAPI.Host)
		},
		ExpectedStatus: http.StatusOK,
	}.Do()
	if err != nil {
		return nil, err
	}
	// enhance response
	ev := response.Event
	GameStatus := ev.Status.Description
	if ev.Status.Type == BasketAPI.StatusType_Inprogress && ev.Time.Played != nil {
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

func (impl *VersionedImpl) RegisterGetGameDetails(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID: "get-game-details",
				Summary:     "Get game details",
				Description: "Return details/stats for a single game",
				Method:      http.MethodGet,
				Errors: []int{
					http.StatusUnauthorized,
					http.StatusBadRequest,
					http.StatusFailedDependency,
				},
				Tags: []string{"BasketAPI"},
				Path: "/game",
			},
		),
		func(ctx context.Context, input *GetGameDetailsInput) (*GetGameDetailsOutput, error) {
			teamEnhancer, err := BasketAPI.GetTeamEnhancer(ctx)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnknownError,
					errors.New("unable to initialize team entity enhancer"),
					err,
				)
			}
			matchDetails, err := getMatchDetails(input.GameId)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err424_BasketAPIGetGameDetails,
					errors.New("unable to retrieve match details"),
					err,
				)
			}
			playersStats, err := getPlayersStats(input.GameId)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err424_BasketAPIGetGameDetails,
					errors.New("unable to retrieve players stats"),
					err,
				)
			}
			teamsStats, err := getTeamStats(input.GameId)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err424_BasketAPIGetGameDetails,
					errors.New("unable to retrieve teams stats"),
					err,
				)
			}
			result := GameDetails{
				MatchDetails: *matchDetails,
				HomeTeam: TeamInfoExtended{
					TeamInfo: teamEnhancer(matchDetails.Event.HomeTeam),
					Players:  playersStats.HomeTeam,
				},
				AwayTeam: TeamInfoExtended{
					TeamInfo: teamEnhancer(matchDetails.Event.AwayTeam),
					Players:  playersStats.AwayTeam,
				},
			}
			if teamsStats != nil {
				result.HomeTeam.Stats = &teamsStats.HomeTeam
				result.AwayTeam.Stats = &teamsStats.AwayTeam
			}
			return &GetGameDetailsOutput{
				Body: result,
			}, nil
		},
	)
}
