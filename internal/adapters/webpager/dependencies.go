package webpager

import (
	"context"

	"demoscraper/internal/core/entities"
)

type HTTPClient interface {
	Get(context.Context, entities.Request) (entities.Response, error)
}
