package misc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type GetList[T any] struct {
	Client         http.Client
	URLs           []string
	Concurrency    int
	RPS            int
	ExpectedStatus int
	UpdateRequest  func(req *http.Request)
}

func (g GetList[T]) Do() ([]T, error) {
	results := []T{}
	source := g.Produce()
	if g.RPS > 0 {
		source = g.Throttle(source)
	}
	for res := range g.SplitAndRun(source) {
		if g.ExpectedStatus != 0 && res.StatusCode != g.ExpectedStatus {
			log.Printf("GetList: one of the requests (%s) failed: %s", res.Request.URL, res.Status)
			continue
		}
		body := res.Body
		defer body.Close()
		var dataItem T
		if err := json.NewDecoder(body).Decode(&dataItem); err != nil {
			return nil, fmt.Errorf("GetList: unable to decode response: %w", err)
		}
		results = append(results, dataItem)
	}
	if len(results) < len(g.URLs) {
		return nil, fmt.Errorf("GetList: result list is too short")
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

func (g GetList[T]) Throttle(in <-chan *http.Request) chan *http.Request {
	ch := make(chan *http.Request)
	ticker := time.NewTicker(time.Second / time.Duration(g.RPS*g.Concurrency))
	go func() {
		for id := range in {
			<-ticker.C
			ch <- id
		}
		ticker.Stop()
		close(ch)
	}()
	return ch
}

func (g GetList[T]) SplitAndRun(in <-chan *http.Request) chan *http.Response {
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
				for req := range in {
					if res, err := g.Client.Do(req); err == nil {
						ch <- res
					} else {
						log.Printf("unable to execute request: %s", err)
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
