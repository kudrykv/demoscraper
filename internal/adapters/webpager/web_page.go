package webpager

import (
	"bytes"
	"context"
	"demoscraper/internal/core/entities"
	"fmt"
	"io"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type WebPage struct {
	httpClient HTTPClient

	rawURL string

	node *html.Node
}

func NewWebPage(httpClient HTTPClient, rawURL string) *WebPage {
	return &WebPage{
		httpClient: httpClient,
		rawURL:     rawURL,
	}
}

func (r *WebPage) Load(ctx context.Context) error {
	response, err := r.httpClient.Get(ctx, entities.Request{URL: r.rawURL})
	if err != nil {
		return fmt.Errorf("get request: %w", err)
	}

	if r.node, err = htmlquery.Parse(io.NopCloser(bytes.NewReader(response.Body))); err != nil {
		return fmt.Errorf("parse html: %w", err)
	}

	return nil
}

func (r *WebPage) Links(_ context.Context) (entities.Links, error) {
	if r.node == nil {
		return nil, ErrNotLoaded
	}

	nodesWithHrefs := htmlquery.Find(r.node, "//*[@href]")
	links := make(entities.Links, 0, len(nodesWithHrefs))

	for _, node := range nodesWithHrefs {
		href := htmlquery.SelectAttr(node, "href")

		link, err := entities.NewLinkFromRawURL(href)
		if err != nil {
			continue
		}

		links = append(links, link)
	}

	return links, nil
}

func (r *WebPage) URL() string {
	return r.rawURL
}

type WebPages []*WebPage
