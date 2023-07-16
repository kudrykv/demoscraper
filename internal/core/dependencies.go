package core

import (
	"context"
	"demoscraper/internal/core/entities"
)

type HTTPClient interface {
	Get(request entities.Request) (response entities.Response, err error)
}

type WebPager interface {
	New(string) WebPage
}

type WebPage interface {
	Load(context.Context) error
}

type WebPages []WebPage
