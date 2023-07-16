package webpager

import (
	"bytes"
	"context"
	"demoscraper/internal/core/entities"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io"
)

type WebPage struct {
	httpClient HTTPClient

	rawURL string
	node   *html.Node
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

	if r.node, err = htmlquery.Parse(io.NopCloser(bytes.NewReader(response.Body))); err != nil {
		return fmt.Errorf("parse html: %w", err)
	}

	return nil
}

type WebPages []*WebPage
