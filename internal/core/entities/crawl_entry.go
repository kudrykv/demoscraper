package entities

type CrawlEntry struct {
	Link  Link
	Depth int
}

func (r CrawlEntry) URL() string {
	return r.Link.URL()
}

func (r CrawlEntry) MatchesHostname(hostname string) bool {
	return r.Link.Hostname() == hostname
}

func NewCrawlEntry(link Link, depth int) CrawlEntry {
	return CrawlEntry{
		Link:  link,
		Depth: depth,
	}
}

type CrawlEntries []CrawlEntry
