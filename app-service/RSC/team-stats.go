package RSC

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func (client *Client) GetTeamStats(query url.Values) ([]TeamSeasonStatItem, error) {
	for queryKey := range query {
		client.Query.Add(queryKey, query.Get(queryKey))
	}
	season := client.GetSeason()
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/team-stats/%s/%s?%s", client.URL, season, client.Sport, client.Query.Encode()),
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
		return []TeamSeasonStatItem{}, nil
	}
	body := res.Body
	defer body.Close()
	var data TeamSeasonStats
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, fmt.Errorf("unable to parse response from the RSC request: %w", err)
	}

	return data.Data.NBA, nil
}
