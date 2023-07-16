package core

import (
	"context"
	"demoscraper/internal/core/entities"
	"fmt"
	"net/url"
)

type Crawler struct{}

func NewCrawler() Crawler {
	return Crawler{}
}

type CrawlParameters struct {
	StartURL string
}

func (r Crawler) Crawl(ctx context.Context, parameters CrawlParameters) (<-chan entities.CrawlEntry, error) {
	_, err := url.Parse(parameters.StartURL)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	entries := make(chan entities.CrawlEntry, 1)
	go func() {
		defer close(entries)

	}()

	return entries, nil
}
