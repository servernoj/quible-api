package RSC

import (
	"encoding/json"
	"fmt"
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

func NewClient[T ResponseItem]() *Client[T] {
	query := url.Values{}
	query.Add("RSC_token", os.Getenv("ENV_RSC_TOKEN"))
	return &Client[T]{
		Client: *http.DefaultClient,
		URL:    BASE_URL,
		Sport:  DEFAULT_SPORT,
		Query:  query,
	}
}

type Client[T ResponseItem] struct {
	http.Client
	URL   string
	Sport string
	Query url.Values
}

func (client *Client[T]) GetDate() string {
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

func (client *Client[T]) GetSeason() string {
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

func (client *Client[T]) RequestRunner(url string) ([]T, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		url,
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create RSC request: %w", err)
	}
	log.Println("making RSC request to", req.URL.String())
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to execute the RSC request: %w", err)
	}
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("RSC request returned error: %s", res.Status)
	}
	if res.StatusCode == 304 {
		return []T{}, nil
	}
	body := res.Body
	defer body.Close()
	var data Response[T]
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, fmt.Errorf("unable to parse response from the RSC request: %w", err)
	}
	return data.Data.NBA, nil
}
