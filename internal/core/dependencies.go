package core

import (
	"context"
	"demoscraper/internal/core/entities"
)

type WebPager interface {
	New(rawURL string) WebPage
	NewFromLinks(entities.Links) WebPages
}

type WebPage interface {
	Load(context.Context) error
	Links(context.Context) (entities.Links, error)
	URL() string
}

type WebPages []WebPage
