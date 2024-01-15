package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GetOne[T any] struct {
	Client        http.Client
	URL           string
	UpdateRequest func(req *http.Request)
}

func (g GetOne[T]) Do() (*T, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		g.URL,
		http.NoBody,
	)
	if g.UpdateRequest != nil {
		g.UpdateRequest(request)
	}
	if err != nil {
		return nil, fmt.Errorf("GetOne: unable to prepare request: %w", err)
	}
	res, err := g.Client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("GetOne: unable to execute the request: %w", err)
	}
	body := res.Body
	defer body.Close()
	var dataItem T
	if err := json.NewDecoder(body).Decode(&dataItem); err != nil {
		return nil, fmt.Errorf("GetOne: unable to decode response: %w", err)
	}
	return &dataItem, nil
}
