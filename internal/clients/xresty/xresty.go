package xresty

import (
	"context"
	"demoscraper/internal/core/entities"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type Client struct {
	client *resty.Client
}

func New() Client {
	return Client{
		client: resty.NewWithClient(http.DefaultClient),
	}
}

func (r Client) Get(ctx context.Context, request entities.Request) (entities.Response, error) {
	response, err := r.client.R().SetContext(ctx).Get(request.URL)
	if err != nil {
		return entities.Response{}, fmt.Errorf("get request: %w", err)
	}

	return entities.Response{
		Body: response.Body(),
	}, nil
}
