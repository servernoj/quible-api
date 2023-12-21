package main

import (
	"log"
	"os"

	"github.com/quible-io/quible-api/cmd/crawl/espn"
	"github.com/quible-io/quible-api/lib/env"
	"github.com/quible-io/quible-api/lib/store"
)

type Crawler interface {
	Run() error
}

func main() {
	// -- Environment vars from .env file
	env.Setup()
	// -- Store + ORM
	if err := store.Setup(os.Getenv("ENV_DSN")); err != nil {
		log.Fatalf("unable to setup DB connection: %s", err)
	}
	defer store.Close()
	// -- Setup crawler
	var crawler Crawler
	crawlerImpl := "espn"
	if len(os.Args[1]) != 0 {
		crawlerImpl = os.Args[1]
	}
	switch crawlerImpl {
	default:
		crawler = espn.NewCrawler()
	}
	// -- Run crawler
	if err := crawler.Run(); err != nil {
		log.Fatalf("unable to run crawler %q: %s", crawlerImpl, err)
	}
}
