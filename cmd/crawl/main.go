package main

import (
	"fmt"
	"os"

	"github.com/quible-io/quible-api/cmd/crawl/espn"
)

type Crawler interface {
	GetURL() string
}

func main() {
	var crawler Crawler
	crawlerImpl := "espn"
	if len(os.Args[1]) != 0 {
		crawlerImpl = os.Args[1]
	}

	switch crawlerImpl {
	case "espn":
		crawler = espn.NewCrawler()
	}

	fmt.Println("ESPN", crawler.GetURL())
}
