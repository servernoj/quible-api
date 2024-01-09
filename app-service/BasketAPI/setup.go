package BasketAPI

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const HOST = "basketapi1.p.rapidapi.com"
const ERRORS_IN_A_ROW_TO_SET_ALERT = 10
const OK_IN_A_ROW_TO_CLEAR_ALERT = 10

func Setup() chan<- struct{} {
	quit := make(chan struct{})
	ticker := time.NewTicker(2 * time.Second)
	countError := uint(0)
	countOK := uint(0)
	isInError := false
	go func() {
		for {
			select {
			case <-ticker.C:
				url := fmt.Sprintf("https://%s/api/basketball/matches/live", HOST)
				req, _ := http.NewRequest("GET", url, nil)
				req.Header.Add("X-RapidAPI-Key", os.Getenv("ENV_RAPIDAPI_KEY"))
				req.Header.Add("X-RapidAPI-Host", HOST)
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					countError++
					countOK = 0
					if !isInError && countError > ERRORS_IN_A_ROW_TO_SET_ALERT {
						isInError = true
						// TODO: send email alert
					}
					log.Printf("BasketAPI error: %s", err)
				} else {
					countOK++
					countError = 0
					if isInError && countOK > OK_IN_A_ROW_TO_CLEAR_ALERT {
						isInError = false
					}
				}
				var body Live
				if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
					log.Printf(
						"BasketAPI error: %s",
						fmt.Errorf("unable to decode response: %w", err),
					)
				}

				for _, ev := range body.Events {
					if ev.Tournament.Name == "NBA" {
						log.Printf(
							"%s: %d:%d (%s)\n",
							ev.Slug,
							ev.HomeScore.Current,
							ev.AwayScore.Current,
							ev.Status.Description,
						)
					}
				}
				log.Println("")
				res.Body.Close()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return quit
}
