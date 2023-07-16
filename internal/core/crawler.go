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

	processedLinks := entities.Links{link}
	visitedMap := make(map[string]struct{})

	go func() {
		defer close(processedCrawlEntriesChan)

		for depth := 1; depth <= parameters.DepthLimit; depth++ {
			webPages := r.webPager.NewFromLinks(processedLinks)
			processedLinks = processedLinks[:0]

			for _, page := range webPages {
				if err := page.Load(ctx); err != nil {
					continue
				}

				links, err := page.Links(ctx)
				if err != nil {
					continue
				}

				links = links.SupplementMissingHostname(link).FilterHostname(hostname).Cleanup().Unique().DropVisited(visitedMap)
				visitedMap = r.merge(visitedMap, links.ToVisitedMap())

				for _, entry := range links.ToCrawlEntries(depth) {
					processedCrawlEntriesChan <- entry
				}

				processedLinks = append(processedLinks, links...)
			}
		}
	}()

	return processedCrawlEntriesChan, nil
}

func (r *Crawler) merge(left map[string]struct{}, right map[string]struct{}) map[string]struct{} {
	for k, v := range right {
		left[k] = v
	}

	return left
}
