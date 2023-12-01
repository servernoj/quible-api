package RSC

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	BASE_URL      = "http://reliantstats.com/api/v1"
	DEFAULT_SPORT = "NBA"
)

func NewClient() *Client {
	query := url.Values{}
	query.Add("RSC_token", os.Getenv("ENV_RSC_TOKEN"))
	return &Client{
		Client: *http.DefaultClient,
		URL:    BASE_URL,
		Sport:  DEFAULT_SPORT,
		Query:  query,
	}
}

type Client struct {
	http.Client
	URL   string
	Sport string
	Query url.Values
}

func (client *Client) GetDate() string {
	date := "now"
	if client.Query.Has("date") {
		if ts, err := time.Parse(time.DateOnly, client.Query.Get("date")); err == nil {
			date = ts.Format(time.DateOnly)
			client.Query.Del("date")
		} else {
			log.Printf("unable to parse `date` into YYYY-MM-DD format: %s", err)
		}
	}
	return date
}

func (client *Client) GetSeason() string {
	// -- defaults to "current season"
	date := ""
	dateLayout := "2006"
	if client.Query.Has("date") {
		if ts, err := time.Parse(dateLayout, client.Query.Get("date")); err == nil {
			date = ts.Format(dateLayout)
			client.Query.Del("date")
		} else {
			log.Printf("unable to parse `date` into YYYY format: %s", err)
		}
	}
	return date
}
