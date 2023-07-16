package core

import (
	"context"
	"demoscraper/internal/core/entities"
)

type Crawler struct {
	webPager WebPager
}

func NewCrawler(webPager WebPager) Crawler {
	return Crawler{
		webPager: webPager,
	}
}

type CrawlParameters struct {
	StartURL string
}

func (r Crawler) Crawl(ctx context.Context, parameters CrawlParameters) (<-chan entities.CrawlEntry, error) {
	entries := make(chan entities.CrawlEntry, 1)
	pagesToVisit := WebPages{r.webPager.New(parameters.StartURL)}

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

					webPage.Links()
				}
			}
		}
	}()

	return entries, nil
}
