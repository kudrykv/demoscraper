package xresty

import (
	"demoscraper/internal/core/entities"
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

func (r Client) Get(_ entities.Request) (entities.Response, error) {
	return entities.Response{}, nil
}
