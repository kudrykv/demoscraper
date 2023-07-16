package core

import (
	"context"
	"demoscraper/internal/core/entities"
	"fmt"
)

type Crawler struct {
	webPager WebPager
}

func NewCrawler(webPager WebPager) *Crawler {
	return &Crawler{
		webPager: webPager,
	}
}

type CrawlParameters struct {
	StartURL   string
	DepthLimit int
}

func (r *Crawler) Crawl(ctx context.Context, parameters CrawlParameters) (<-chan entities.CrawlEntry, error) {
	link, err := entities.NewLinkFromRawURL(parameters.StartURL)
	if err != nil {
		return nil, fmt.Errorf("link from raw url: %w", err)
	}

	hostname := link.Hostname()

	processedCrawlEntriesChan := make(chan entities.CrawlEntry, 1)

	pagesToVisitChan := make(chan WebPage, 1)
	crawlEntriesChan := make(chan entities.CrawlEntry, 1)
	//visitedPages := make(map[string]struct{})

	pagesToVisitChan <- r.webPager.New(parameters.StartURL, 0)

	go func() {
		defer close(crawlEntriesChan)

		for webPage := range pagesToVisitChan {
			if webPage.Depth() > parameters.DepthLimit {
				continue
			}

			if err := webPage.Load(ctx); err != nil {
				continue
			}

			links, err := webPage.Links(ctx)
			if err != nil {
				continue
			}

			for _, entry := range links.ToCrawlEntries(webPage.Depth() + 1) {
				crawlEntriesChan <- entry
			}
		}
	}()

	go func() {
		defer close(processedCrawlEntriesChan)

		for crawlEntry := range crawlEntriesChan {
			if !crawlEntry.MatchesHostname(hostname) {
				continue
			}

			processedCrawlEntriesChan <- crawlEntry

			if crawlEntry.Depth < parameters.DepthLimit {
				pagesToVisitChan <- r.webPager.New(crawlEntry.URL(), crawlEntry.Depth)
			}
		}
	}()

	return processedCrawlEntriesChan, nil
}
