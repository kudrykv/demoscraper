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
	StartURL string
}

func (r *Crawler) Crawl(ctx context.Context, parameters CrawlParameters) (<-chan entities.CrawlEntry, error) {
	link, err := entities.NewLinkFromRawURL(parameters.StartURL)
	if err != nil {
		return nil, fmt.Errorf("link from raw url: %w", err)
	}

	entries := make(chan entities.CrawlEntry, 1)
	pagesToVisit := WebPages{r.webPager.New(parameters.StartURL)}

	go func() {
		defer close(entries)

		for depth := 1; ; depth++ {
			select {
			case <-ctx.Done():
				return

			default:
				for _, webPage := range pagesToVisit {
					if err := webPage.Load(ctx); err != nil {
						continue
					}

					links, err := webPage.Links(ctx)
					if err != nil {
						continue
					}

					links = links.Unique().FilterHostname(link.Hostname())

					_ = links
				}
			}
		}
	}()

	return entries, nil
}
