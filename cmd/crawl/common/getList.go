package common

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
)

type GetList[T any] struct {
	Client        http.Client
	URLs          []string
	Concurrency   int
	UpdateRequest func(req *http.Request)
}

func (g GetList[T]) Do() ([]T, error) {
	results := []T{}
	source := g.Produce()
	for res := range g.SplitAndRun(source) {
		body := res.Body
		defer body.Close()
		var dataItem T
		if err := json.NewDecoder(body).Decode(&dataItem); err != nil {
			return []T{}, fmt.Errorf("GetList: unable to decode response: %w", err)
		}
		results = append(results, dataItem)
	}
	if len(results) < len(g.URLs) {
		log.Println(results)
		return []T{}, fmt.Errorf("GetList: result list is too short")
	}
	return results, nil
}

func (g GetList[T]) Produce() chan *http.Request {
	ch := make(chan *http.Request)
	go func() {
		for _, url := range g.URLs {
			request, err := http.NewRequest(
				http.MethodGet,
				url,
				http.NoBody,
			)
			if err == nil {
				if g.UpdateRequest != nil {
					g.UpdateRequest(request)
				}
				ch <- request
			}
		}
		close(ch)
	}()
	return ch
}

func (g GetList[T]) SplitAndRun(requests <-chan *http.Request) chan *http.Response {
	ch := make(chan *http.Response)
	go func() {
		concurrency := g.Concurrency
		if concurrency == 0 {
			concurrency = runtime.NumCPU()
		}
		var wg sync.WaitGroup
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func() {
				for req := range requests {
					if res, err := g.Client.Do(req); err == nil {
						ch <- res
					}
				}
				wg.Done()
			}()
		}
		wg.Wait()
		close(ch)
	}()
	return ch
}
