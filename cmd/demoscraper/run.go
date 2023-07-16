package main

import (
	"context"
	"demoscraper/internal/adapters/webpager"
	"demoscraper/internal/clients/xresty"
	"demoscraper/internal/core"
	"log"
)

func run(ctx context.Context) {
	httpClient := xresty.New()
	webPager := webpager.New(httpClient)
	crawler := core.NewCrawler(webPager)

	crawlEntries, err := crawler.Crawl(ctx, core.CrawlParameters{
		StartURL:   flagStartingURL,
		DepthLimit: flagDepth,
	})
	if err != nil {
		log.Println(err)

		return
	}

	for crawlEntry := range crawlEntries {
		log.Printf("CrawlEntry: %v\n", crawlEntry)
	}
}
