package core

import (
	"context"
	"demoscraper/internal/core/entities"
)

type WebPager interface {
	New(rawURL string, depth int) WebPage
}

type WebPage interface {
	Load(context.Context) error
	Links(context.Context) (entities.Links, error)
	URL() string
	Depth() int
}

type WebPages []WebPage
