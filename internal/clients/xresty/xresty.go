package xresty

import (
	"context"
	"demoscraper/internal/core/entities"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
}

func New(httpClient *http.Client) Client {
	return Client{
		client: resty.NewWithClient(httpClient),
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
