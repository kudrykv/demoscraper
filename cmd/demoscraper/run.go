package main

import (
	"context"
	"demoscraper/internal/core"
	"log"
)

func run(ctx context.Context) {
	crawlEntries, err := core.NewCrawler().Crawl(ctx, core.CrawlParameters{
		StartURL: flagStartingURL,
	})
	if err != nil {
		log.Println(err)

		return
	}

	_ = crawlEntries
}
