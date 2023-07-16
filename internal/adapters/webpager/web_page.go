package webpager

import (
	"context"
	"demoscraper/internal/core/entities"
	"fmt"
)

type WebPage struct {
	httpClient HTTPClient

	rawURL string
}

func NewWebPage(httpClient HTTPClient, rawURL string) *WebPage {
	return &WebPage{
		httpClient: httpClient,
		rawURL:     rawURL,
	}
}

func (r *WebPage) Load(_ context.Context) error {
	response, err := r.httpClient.Get(entities.Request{
		URL: r.rawURL,
	})
	if err != nil {
		return fmt.Errorf("get request: %w", err)
	}

	_ = response

	return nil
}

type WebPages []*WebPage
