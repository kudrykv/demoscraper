package main

import (
	"context"
	"demoscraper/internal/adapters/tsvmarshaller"
	"demoscraper/internal/adapters/webpager"
	"demoscraper/internal/clients/xresty"
	"demoscraper/internal/core"
	"log"
)

func run(ctx context.Context) {
	httpClient := xresty.New()
	tsvMarshaller := tsvmarshaller.New()
	webPager := webpager.New(httpClient)
	crawler := core.NewCrawler(webPager)
	store := core.NewStore(tsvMarshaller)

	crawlEntries, err := crawler.Crawl(ctx, core.CrawlParameters{
		StartURL:   flagStartingURL,
		DepthLimit: flagDepth,
	})
	if err != nil {
		log.Println(err)

		return
	}

	if err := store.Save(ctx, flagOutputFile, crawlEntries); err != nil {
		log.Println(err)

		return
	}
}
