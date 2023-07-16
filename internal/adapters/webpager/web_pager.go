package webpager

import (
	"demoscraper/internal/core"
)

type WebPager struct {
	httpClient HTTPClient
}

func New(httpClient HTTPClient) WebPager {
	return WebPager{
		httpClient: httpClient,
	}
}

func (r WebPager) New(rawURL string, depth int) core.WebPage {
	return NewWebPage(r.httpClient, rawURL, depth)
}
