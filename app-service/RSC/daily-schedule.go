package RSC

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func (client *Client) GetDailySchedule(query url.Values) ([]ScheduleItem, error) {
	for queryKey := range query {
		client.Query.Add(queryKey, query.Get(queryKey))
	}
	date := client.GetDate()
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/schedule/%s/%s?%s", client.URL, date, client.Sport, client.Query.Encode()),
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
		return []ScheduleItem{}, nil
	}
	body := res.Body
	defer body.Close()
	var data Schedule
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, fmt.Errorf("unable to parse response from the RSC request: %w", err)
	}

	return data.Data.NBA, nil
}
