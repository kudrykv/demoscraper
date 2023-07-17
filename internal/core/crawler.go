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
	StartURL    string
	DepthLimit  int
	Parallelism int
}

func (r *Crawler) Crawl(ctx context.Context, parameters CrawlParameters) (<-chan entities.CrawlEntry, error) {
	link, err := entities.NewLinkFromRawURL(parameters.StartURL)
	if err != nil {
		return nil, fmt.Errorf("link from raw url: %w", err)
	}

	crawlEntriesChan := make(chan entities.CrawlEntry, 1)

	go func() {
		defer close(crawlEntriesChan)

		r.crawl(ctx, parameters, link, crawlEntriesChan)
	}()

	return crawlEntriesChan, nil
}

func (r *Crawler) crawl(
	ctx context.Context,
	parameters CrawlParameters,
	root entities.Link,
	crawlEntriesChan chan<- entities.CrawlEntry,
) {
	processedLinks := entities.Links{root}
	visitor := r.makeVisitor()
	mutex := sync.Mutex{}

	for depth := 1; depth <= parameters.DepthLimit; depth++ {
		webPages := r.webPager.NewFromLinks(processedLinks)
		processedLinks = processedLinks[:0]

		waitGroup, semaphore := makeWaitGroupAndSemaphore(len(webPages), parameters.Parallelism)

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

				links = links.
					SupplementMissingHostname(root).
					FilterHostname(root.Hostname()).
					Cleanup().
					Unique().
					DropVisited(visitor.ToVisitMap())

				visitor.Merge(links.ToVisitMap())

				for _, entry := range links.ToCrawlEntries(depth) {
					crawlEntriesChan <- entry
				}

				mutex.Lock()
				processedLinks = append(processedLinks, links...)
				mutex.Unlock()
			}()
		}

		waitGroup.Wait()
	}
}

func makeWaitGroupAndSemaphore(waitGroupSize int, semaphoreSize int) (*sync.WaitGroup, chan struct{}) {
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(waitGroupSize)

	if semaphoreSize < 1 {
		if numCPU := runtime.NumCPU(); numCPU < 1 {
			semaphoreSize = 1
		} else {
			semaphoreSize = numCPU
		}
	}

	return waitGroup, make(chan struct{}, semaphoreSize)
}
