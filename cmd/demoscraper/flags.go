package main

import "flag"

var (
	flagDepth       int
	flagStartingURL string
	flagOutputFile  string
)

func setupFlags() {
	flag.IntVar(&flagDepth, "depth", 2, "Depth of scraping")
	flag.StringVar(&flagStartingURL, "url", "https://github.com", "URL to start scraping from")
	flag.StringVar(&flagOutputFile, "output", "output.tsv", "Output file")

	flag.Parse()
}
