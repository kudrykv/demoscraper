package xresty

import (
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

func (r Client) Get(request entities.Request) (entities.Response, error) {
	response, err := r.client.R().Get(request.URL)
	if err != nil {
		return entities.Response{}, fmt.Errorf("get request: %w", err)
	}

	return entities.Response{
		Body: response.Body(),
	}, nil
}
