package main

import "flag"

var (
	flagDepth       int
	flagStartingURL string
)

func setupFlags() {
	flag.IntVar(&flagDepth, "depth", 1, "Depth of scraping")
	flag.StringVar(&flagStartingURL, "url", "https://www.google.com", "URL to start scraping from")

	flag.Parse()
}
