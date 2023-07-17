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

func (r Link) URL() string {
	return r.url.String()
}

func (r Link) String() string {
	if r.url == nil {
		return "empty link"
	}

	return r.URL()
}

type Links []Link

func (r Links) Unique() Links {
	if len(r) == 0 {
		return nil
	}

	unique := make(map[string]struct{})
	result := make(Links, 0, len(r))

	for _, link := range r {
		key := link.URL()

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

func (r Links) ToCrawlEntries(depth int) CrawlEntries {
	if len(r) == 0 {
		return nil
	}

	crawlEntries := make(CrawlEntries, 0, len(r))

	for _, link := range r {
		crawlEntries = append(crawlEntries, NewCrawlEntry(link, depth))
	}

	return crawlEntries
}

func (r Links) filter(filterFunc func(link Link) bool) Links {
	if len(r) == 0 {
		return nil
	}

	var result Links

	for _, link := range r {
		if filterFunc(link) {
			result = append(result, link)
		}
	}

	return result
}

func (r Links) SupplementMissingHostname(link Link) Links {
	if len(r) == 0 {
		return nil
	}

	links := make(Links, len(r))
	copy(links, r)

	for i := range links {
		if links[i].Hostname() == "" {
			links[i].url.Scheme = link.url.Scheme
			links[i].url.Host = link.url.Host
		}
	}

	return links
}

func (r Links) Cleanup() Links {
	if len(r) == 0 {
		return nil
	}

	links := make(Links, len(r))
	copy(links, r)

	for i := range links {
		links[i].url.Fragment = ""
		links[i].url.Opaque = ""

		if links[i].url.Path == "/" {
			links[i].url.Path = ""
		}
	}

	return links
}

func (r Links) DropVisited(hitMap VisitMap) Links {
	if len(r) == 0 {
		return nil
	}

	result := make(Links, 0, len(r))

	for _, link := range r {
		if _, ok := hitMap[link.URL()]; ok {
			continue
		}

		result = append(result, link)
	}

	return result
}

func (r Links) ToVisitMap() VisitMap {
	if len(r) == 0 {
		return nil
	}

	visitedMap := make(VisitMap, len(r))
	for _, link := range r {
		visitedMap[link.URL()] = struct{}{}
	}

	return visitedMap
}
