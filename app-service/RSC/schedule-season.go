package RSC

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (client *Client) GetScheduleSeason() ([]ScheduleItem, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/schedule-season/%s?%s", client.URL, client.Sport, client.Query.Encode()),
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
