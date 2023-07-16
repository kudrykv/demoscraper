package core

import (
	"context"
	"demoscraper/internal/core/entities"
	"fmt"
	"sync/atomic"
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

	webPagesNumber := atomic.Int64{}
	crawlNumber := atomic.Int64{}

	webPagesNumber.Store(1)
	crawlNumber.Store(0)

	//visitedPages := make(map[string]struct{})

	pagesToVisitChan <- r.webPager.New(parameters.StartURL, 0)

	go func() {
		defer close(crawlEntriesChan)

		for webPage := range pagesToVisitChan {
			if webPage.Depth() > parameters.DepthLimit {
				if webPagesNumber.Add(-1) == 0 {
					if crawlNumber.Load() == 0 {
						return
					}
				}

				continue
			}

			if err := webPage.Load(ctx); err != nil {
				if webPagesNumber.Add(-1) == 0 {
					if crawlNumber.Load() == 0 {
						return
					}
				}

				continue
			}

			links, err := webPage.Links(ctx)
			if err != nil {
				if webPagesNumber.Add(-1) == 0 {
					if crawlNumber.Load() == 0 {
						return
					}
				}

				continue
			}

			toCrawlEntries := links.ToCrawlEntries(webPage.Depth() + 1)

			crawlNumber.Add(int64(len(toCrawlEntries)))

			for _, entry := range toCrawlEntries {
				crawlEntriesChan <- entry
			}

			if webPagesNumber.Add(-1) == 0 {
				if crawlNumber.Load() == 0 {
					return
				}
			}
		}
	}()

	go func() {
		defer close(processedCrawlEntriesChan)

		for crawlEntry := range crawlEntriesChan {
			if !crawlEntry.MatchesHostname(hostname) {
				if crawlNumber.Add(-1) == 0 {
					if webPagesNumber.Load() == 0 {
						return
					}
				}

				continue
			}

			processedCrawlEntriesChan <- crawlEntry

			if crawlEntry.Depth <= parameters.DepthLimit {
				pagesToVisitChan <- r.webPager.New(crawlEntry.URL(), crawlEntry.Depth)
			}

			if crawlNumber.Add(-1) == 0 {
				if webPagesNumber.Load() == 0 {
					return
				}
			}
		}
	}()

	return processedCrawlEntriesChan, nil
}
