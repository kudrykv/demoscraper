package entities

import (
	"fmt"
	"net/url"
)

type Link struct {
	url *url.URL
}

func NewLinkFromRawURL(rawURL string) (Link, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return Link{}, fmt.Errorf("parse url: %w", err)
	}

	return Link{url: parsed}, nil
}

func (r Link) Hostname() string {
	return r.url.Hostname()
}

type Links []Link

func (r Links) Unique() Links {
	if len(r) == 0 {
		return nil
	}

	unique := make(map[string]struct{})
	var result Links
	for _, link := range r {
		key := link.url.String()

		if _, ok := unique[key]; ok {
			continue
		}

		unique[key] = struct{}{}
		result = append(result, link)
	}

	return result
}

func (r Links) FilterHostname(hostname string) Links {
	return r.filter(func(link Link) bool { return link.Hostname() == hostname })
}

func (r Links) filter(f func(link Link) bool) Links {
	if len(r) == 0 {
		return nil
	}

	var result Links
	for _, link := range r {
		if f(link) {
			result = append(result, link)
		}
	}

	return result
}
