package core

import (
	"demoscraper/internal/core/entities"
	"fmt"
	"net/url"
)

type Crawler struct{}

func NewCrawler() Crawler {
	return Crawler{}
}

func (r Crawler) Crawl(startingURI string) (<-chan entities.CrawlEntry, error) {
	_, err := url.Parse(startingURI)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	entries := make(chan entities.CrawlEntry, 1)
	close(entries)

	return entries, nil
}
