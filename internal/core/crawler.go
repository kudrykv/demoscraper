package core

import (
	"context"
	"demoscraper/internal/core/entities"
	"fmt"
	"runtime"
	"sync"
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
	visitedMap := make(map[string]struct{})

	for depth := 1; depth <= parameters.DepthLimit; depth++ {
		webPages := r.webPager.NewFromLinks(processedLinks)
		processedLinks = processedLinks[:0]

		waitGroup := sync.WaitGroup{}
		waitGroup.Add(len(webPages))

		semaphoreLimit := runtime.NumCPU()
		if semaphoreLimit < 1 {
			semaphoreLimit = 1
		}

		semaphoreLimit = 1

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

				links = links.SupplementMissingHostname(link).FilterHostname(hostname).Cleanup().Unique().DropVisited(visitedMap)
				visitedMap = r.merge(visitedMap, links.ToVisitedMap())

				for _, entry := range links.ToCrawlEntries(depth) {
					processedCrawlEntriesChan <- entry
				}

				processedLinks = append(processedLinks, links...)
			}()
		}

		waitGroup.Wait()
	}
}

func (r *Crawler) processOne(ctx context.Context, root entities.Link, webPage WebPage, visitedMap map[string]struct{}) (entities.Links, error) {
	links, err := r.processWebpage(ctx, webPage)
	if err != nil {
		return nil, fmt.Errorf("process webpage: %w", err)
	}

	links, updatedVisitMap := r.processLinksBatch(links, root, visitedMap)

	ptr := &visitedMap
	*ptr = updatedVisitMap

	panic("later")
}

func (r *Crawler) processWebpage(ctx context.Context, webPage WebPage) (entities.Links, error) {
	if err := webPage.Load(ctx); err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	links, err := webPage.Links(ctx)
	if err != nil {
		return nil, fmt.Errorf("links: %w", err)
	}

	return links, nil
}

func (r *Crawler) processLinksBatch(
	links entities.Links, root entities.Link, visitedMap map[string]struct{},
) (entities.Links, map[string]struct{}) {
	filteredLinks := links.
		SupplementMissingHostname(root).
		FilterHostname(root.Hostname()).
		Cleanup().
		Unique().
		DropVisited(visitedMap)

	return filteredLinks, r.merge(visitedMap, filteredLinks.ToVisitedMap())
}

func (r *Crawler) merge(left map[string]struct{}, right map[string]struct{}) map[string]struct{} {
	for k, v := range right {
		left[k] = v
	}

	return left
}
