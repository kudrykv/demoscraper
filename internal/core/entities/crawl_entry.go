package entities

import "fmt"

type CrawlEntry struct {
	Link  Link
	Depth int
}

func NewCrawlEntry(link Link, depth int) CrawlEntry {
	return CrawlEntry{
		Link:  link,
		Depth: depth,
	}
}

func (r CrawlEntry) URL() string {
	return r.Link.URL()
}

func (r CrawlEntry) MatchesHostname(hostname string) bool {
	return r.Link.Hostname() == hostname
}

func (r CrawlEntry) String() string {
	return fmt.Sprintf("%+v at %d", r.Link, r.Depth)
}

type CrawlEntries []CrawlEntry
