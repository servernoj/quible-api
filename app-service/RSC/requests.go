package RSC

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
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

func (client *Client) GetScheduleSeason() ([]ScheduleSeasonItem, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/schedule-season/%s??%s", client.URL, client.Sport, client.Query.Encode()),
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create RSC request: %w", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to execute the RSC request: %w", err)
	}
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("RSC request returned error: %s", res.Status)
	}
	body := res.Body
	defer body.Close()
	var data ScheduleSeason
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, fmt.Errorf("unable to parse response from the RSC request: %w", err)
	}

	return data.Data.NBA, nil
}
