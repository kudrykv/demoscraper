package core

import (
	"context"
)

type WebPage struct {
	rawURL string
}

func NewWebPage(rawURL string) *WebPage {
	return &WebPage{
		rawURL: rawURL,
	}
}

func (r *WebPage) Load(_ context.Context) error {
	panic("not implemented")
}

type WebPages []*WebPage
