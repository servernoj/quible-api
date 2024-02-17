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

type ListGamesInput struct {
	AuthorizationHeaderResolver
	Date               string `query:"date" format:"date"`
	LocalTimeZoneShift int    `query:"localTimeZoneShift" exclusiveMaximum:"0"`
}

type ListGamesOutput struct {
	Body []Game
}

func (impl *VersionedImpl) RegisterListGames(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID: "get-games",
				Summary:     "Get games",
				Description: "List games scheduled for the given date",
				Method:      http.MethodGet,
				Errors: []int{
					http.StatusUnauthorized,
					http.StatusBadRequest,
					http.StatusFailedDependency,
				},
				Tags: []string{"BasketAPI"},
				Path: "/games",
			},
		),
		func(ctx context.Context, input *ListGamesInput) (*ListGamesOutput, error) {
			// 1. Adjust client location to narrow down list of games
			loc, _ := time.LoadLocation("America/New_York")
			if input.LocalTimeZoneShift < 0 {
				loc = time.FixedZone("User timezone", input.LocalTimeZoneShift*int(time.Hour/time.Second))
			}
			dateParsed, _ := time.Parse(time.DateOnly, input.Date)
			dateParsedInLocation, _ := time.ParseInLocation(time.DateOnly, input.Date, loc)
			// 2. Send request to Matches API
			url := fmt.Sprintf(
				"https://%s/api/basketball/matches/%s",
				BasketAPI.Host,
				dateParsed.Format("2/1/2006"),
			)
			response, err := misc.GetOne[MS_Data]{
				Client: *http.DefaultClient,
				URL:    url,
				UpdateRequest: func(req *http.Request) {
					req.Header.Set("X-RapidAPI-Key", os.Getenv("ENV_RAPIDAPI_KEY"))
					req.Header.Set("X-RapidAPI-Host", BasketAPI.Host)
				},
				ExpectedStatus: http.StatusOK,
			}.Do()
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err424_BasketAPIListGames,
					errors.New("unable to retrieve matches"),
					err,
				)
			}
			// 3. Initialize team data enhancer (inject logo, arena, colors, etc)
			teamEnhancer, err := BasketAPI.GetTeamEnhancer(ctx)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnknownError,
					errors.New("unable to initialize team entity enhancer"),
					err,
				)
			}
			// 4. Convert "events" into "games"
			games := make([]Game, 0, len(response.Events))
			tsFrom := dateParsedInLocation.Unix()
			tsTo := dateParsedInLocation.Add(24 * time.Hour).Unix()
			for _, ev := range response.Events {
				if ev.Tournament.Name != "NBA" || ev.StartTimestamp < tsFrom || ev.StartTimestamp > tsTo {
					continue
				}
				ev := ev
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
			// 5. Send response with the list of games
			return &ListGamesOutput{
				Body: games,
			}, nil
		},
	)
}
