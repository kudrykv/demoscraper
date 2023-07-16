package webpager

import (
	"demoscraper/internal/core"
	"demoscraper/internal/core/entities"
)

type WebPager struct {
	httpClient HTTPClient
}

func New(httpClient HTTPClient) WebPager {
	return WebPager{
		httpClient: httpClient,
	}
}

func (r WebPager) New(rawURL string) core.WebPage {
	return NewWebPage(r.httpClient, rawURL)
}

func (r WebPager) NewFromLinks(links entities.Links) core.WebPages {
	if len(links) == 0 {
		return nil
	}

	webPages := make(core.WebPages, 0, len(links))

	for _, link := range links {
		webPages = append(webPages, r.New(link.URL()))
	}

	return webPages
}
