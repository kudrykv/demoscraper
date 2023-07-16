package webpager

import "demoscraper/internal/core/entities"

type HTTPClient interface {
	Get(request entities.Request) (response entities.Response, err error)
}
