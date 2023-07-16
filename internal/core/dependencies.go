package core

import (
	"context"
	"demoscraper/internal/core/entities"
)

type WebPager interface {
	New(string) WebPage
}

type WebPage interface {
	Load(context.Context) error
	Links(context.Context) (entities.Links, error)
}

type WebPages []WebPage
