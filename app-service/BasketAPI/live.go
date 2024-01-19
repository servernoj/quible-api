package BasketAPI

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/ably/ably-go/ably"
	"github.com/quible-io/quible-api/lib/email"
	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/models"
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
	// ably
	ablyRealTime, err := ably.NewRealtime(
		ably.WithKey(os.Getenv("ENV_ABLY_KEY")),
		ably.WithClientID("backend"),
	)
	if err != nil {
		log.Printf("unable to initialize Ably SDK: %s", err)
		return nil, err
	}
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
						ev.HomeTeam.Logo = logoByShortName[ev.HomeTeam.ShortName]
						ev.AwayTeam.Logo = logoByShortName[ev.AwayTeam.ShortName]
						liveMessage.Events = append(liveMessage.Events, ev)
						states[ev.ID] = state
						log.Printf(
							"[%d] %s: %d:%d (%s)\n",
							ev.ID,
							ev.Slug,
							ev.AwayScore.Current,
							ev.HomeScore.Current,
							ev.Status.Description,
						)
					}
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
