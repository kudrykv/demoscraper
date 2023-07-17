package core

import (
	"context"
	"demoscraper/internal/core/entities"
	"fmt"
	"runtime"
	"sync"
)

type Crawler struct {
	webPager    WebPager
	makeVisitor MakeVisitor
}

func NewCrawler(webPager WebPager, makeVisitor MakeVisitor) *Crawler {
	return &Crawler{
		webPager:    webPager,
		makeVisitor: makeVisitor,
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

	processedCrawlEntriesChan := make(chan entities.CrawlEntry, 1)

	go func() {
		defer close(processedCrawlEntriesChan)

		r.crawl(ctx, parameters, link, processedCrawlEntriesChan)
	}()

	return processedCrawlEntriesChan, nil
}

func (r *Crawler) crawl(
	ctx context.Context,
	parameters CrawlParameters,
	link entities.Link,
	processedCrawlEntriesChan chan<- entities.CrawlEntry,
) {
	hostname := link.Hostname()
	processedLinks := entities.Links{link}
	visitor := r.makeVisitor()
	mutex := sync.Mutex{}

	for depth := 1; depth <= parameters.DepthLimit; depth++ {
		webPages := r.webPager.NewFromLinks(processedLinks)
		processedLinks = processedLinks[:0]

		waitGroup := sync.WaitGroup{}
		waitGroup.Add(len(webPages))

		semaphoreLimit := runtime.NumCPU()
		if semaphoreLimit < 1 {
			semaphoreLimit = 1
		}

		semaphore := make(chan struct{}, semaphoreLimit)

		for _, webPage := range webPages {
			webPage := webPage
			semaphore <- struct{}{}

			go func() {
				defer waitGroup.Done()
				defer func() { <-semaphore }()

				if err := webPage.Load(ctx); err != nil {
					return
				}

				links, err := webPage.Links(ctx)
				if err != nil {
					return
				}

				links = links.SupplementMissingHostname(link).FilterHostname(hostname).Cleanup().Unique()

				links = links.DropVisited(visitor.ToVisitMap())
				visitor.Merge(links.ToVisitedMap())

				for _, entry := range links.ToCrawlEntries(depth) {
					processedCrawlEntriesChan <- entry
				}

				mutex.Lock()
				processedLinks = append(processedLinks, links...)
				mutex.Unlock()
			}()
		}

		waitGroup.Wait()
	}
}
