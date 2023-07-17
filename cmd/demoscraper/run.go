package main

import (
	"context"
	"demoscraper/internal/adapters/inmemvisitor"
	"demoscraper/internal/adapters/tsvmarshaller"
	"demoscraper/internal/adapters/webpager"
	"demoscraper/internal/clients/xresty"
	"demoscraper/internal/core"
	"demoscraper/internal/core/entities"
	"log"
	"net/http"
	"sync"
)

func run(ctx context.Context) {
	httpClient := xresty.New(http.DefaultClient)
	tsvMarshaller := tsvmarshaller.New()
	webPager := webpager.New(httpClient)
	crawler := core.NewCrawler(webPager, inmemvisitor.New)
	store := core.NewStore(tsvMarshaller)

	crawlEntries, err := crawler.Crawl(ctx, core.CrawlParameters{
		StartURL:    flagStartingURL,
		DepthLimit:  flagDepth,
		Parallelism: flagParallelism,
	})
	if err != nil {
		log.Println(err)

		return
	}

	logChan := make(chan entities.CrawlEntry, 1)
	saveChan := make(chan entities.CrawlEntry, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for entry := range crawlEntries {
			logChan <- entry
			saveChan <- entry
		}

		close(logChan)
		close(saveChan)
		wg.Done()
	}()

	go func() {
		for entry := range logChan {
			log.Println(entry)
		}
	}()

	if err := store.Save(ctx, flagOutputFile, saveChan); err != nil {
		log.Println(err)

		return
	}

	wg.Wait()
}
