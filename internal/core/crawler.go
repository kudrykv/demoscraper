package core

import (
	"context"
	"demoscraper/internal/core/entities"
)

type Crawler struct{}

func NewCrawler() Crawler {
	return Crawler{}
}

type CrawlParameters struct {
	StartURL string
}

func (r Crawler) Crawl(ctx context.Context, parameters CrawlParameters) (<-chan entities.CrawlEntry, error) {
	entries := make(chan entities.CrawlEntry, 1)
	pagesToVisit := WebPages{NewWebPage(parameters.StartURL)}

	go func() {
		defer close(entries)

		for {
			select {
			case <-ctx.Done():
				return

			default:
				for _, webPage := range pagesToVisit {
					if err := webPage.Load(ctx); err != nil {
						continue
					}
				}
			}
		}
	}()

	return entries, nil
}
