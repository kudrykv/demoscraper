package main

import "flag"

//nolint:gochecknoglobals
var (
	flagDepth       int
	flagStartingURL string
	flagOutputFile  string
	flagParallelism int
)

func setupFlags() {
	flag.IntVar(&flagDepth, "depth", 2, "Depth of scraping") //nolint:gomnd
	flag.StringVar(&flagStartingURL, "url", "https://github.com", "URL to start scraping from")
	flag.StringVar(&flagOutputFile, "output", "output.tsv", "Output file")
	flag.IntVar(&flagParallelism, "parallelism", 0, "Number of parallel workers (0 = number of CPUs)")

	flag.Parse()
}
