package BasketAPI

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/quible-io/quible-api/app-service-ng/services/ablyService"
	"github.com/quible-io/quible-api/lib/email"
	"github.com/quible-io/quible-api/lib/misc"
)

const ERRORS_IN_A_ROW_TO_SET_ALERT = 10
const OK_IN_A_ROW_TO_CLEAR_ALERT = 10

func StartLive() (chan<- struct{}, error) {
	ctx := context.Background()
	quit := make(chan struct{})
	ticker := time.NewTicker(2 * time.Second)
	countError := uint(0)
	countOK := uint(0)
	isInError := false
	states := map[uint]string{}
	teamEnhancer, err := GetTeamEnhancer(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize team entity enhancer: %w", err)
	}
	// ably
	ablyRealTime := ablyService.GetAbly()
	ablyChannel := ablyRealTime.Channels.Get("live:main")
	go func() {
		for {
			select {
			case <-ticker.C:
				url := fmt.Sprintf("https://%s/api/basketball/matches/live", Host)
				res, err := misc.GetOne[LM_Data]{
					Client: *http.DefaultClient,
					URL:    url,
					UpdateRequest: func(req *http.Request) {
						req.Header.Set("X-RapidAPI-Key", os.Getenv("ENV_RAPIDAPI_KEY"))
						req.Header.Set("X-RapidAPI-Host", Host)
					},
					ExpectedStatus: http.StatusOK,
				}.Do()
				if err != nil {
					countError++
					countOK = 0
					log.Printf("BasketAPI error: %s", err)
					if !isInError && countError > ERRORS_IN_A_ROW_TO_SET_ALERT {
						isInError = true
						if err := email.Send(ctx, email.EmailDTO{
							From:    "api@quible.io",
							To:      "devops@quible.io",
							Subject: "BasketAPI failure",
							TextBody: fmt.Sprintf(
								"At least %d API errors happened in a row\nThe last error reads %q",
								ERRORS_IN_A_ROW_TO_SET_ALERT,
								err,
							),
						}); err != nil {
							log.Printf("unable to send BasketAPI error report: %s", err)
						}
					}
					continue
				} else {
					countOK++
					countError = 0
					if isInError && countOK > OK_IN_A_ROW_TO_CLEAR_ALERT {
						isInError = false
					}
				}
				// -- process and publish to clients
				var liveMessage LiveMessage
				tournaments := []string{"NBA"}
				for _, ev := range res.Events {
					if slices.Index(tournaments, ev.Tournament.Name) != -1 {
						liveMessage.IDs = append(liveMessage.IDs, ev.ID)
						state := fmt.Sprintf("%d:%d@%s+%d", ev.HomeScore.Current, ev.AwayScore.Current, ev.Status.Description, ev.Time.Played)
						value, ok := states[ev.ID]
						if ok && value == state {
							continue
						}
						liveEvent := LiveEvent{
							ID:             ev.ID,
							Status:         ev.Status,
							HomeTeam:       teamEnhancer(ev.HomeTeam),
							AwayTeam:       teamEnhancer(ev.AwayTeam),
							HomeScore:      ev.HomeScore,
							AwayScore:      ev.AwayScore,
							Time:           ev.Time,
							StartTimestamp: ev.StartTimestamp,
						}
						liveMessage.Events = append(liveMessage.Events, liveEvent)
						states[ev.ID] = state
						log.Printf(
							"[%d %s %s %d] played %d: %s (%d) vs. %s (%d)\n",
							ev.ID,
							ev.Status.Description,
							ev.Status.Type,
							ev.Status.Code,
							ev.Time.Played,
							liveEvent.AwayTeam.Abbr,
							liveEvent.AwayScore.Current,
							liveEvent.HomeTeam.Abbr,
							liveEvent.HomeScore.Current,
						)
					}
				}
				if len(liveMessage.IDs) == 0 && len(states) > 0 {
					log.Println("no events in qualified tournaments...")
					states = map[uint]string{}
				}
				if len(liveMessage.Events) > 0 {
					if err := ablyChannel.Publish(ctx, "message", liveMessage); err != nil {
						log.Printf("unable to publish live data to Ably: %s", err)
					}
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit, nil
}
