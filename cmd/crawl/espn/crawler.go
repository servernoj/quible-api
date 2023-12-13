package espn

import (
	"net/http"
)

const (
	BASE_URL = "https://sports.core.api.espn.com/v2/sports/basketball/leagues/nba"
)

func NewCrawler() *Crawler {
	return &Crawler{
		Client: *http.DefaultClient,
		URL:    BASE_URL,
	}
}

type Crawler struct {
	http.Client
	URL string
}

func (c *Crawler) GetURL() string {
	return c.URL
}
