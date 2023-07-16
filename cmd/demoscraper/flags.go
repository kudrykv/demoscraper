package main

import "flag"

var (
	flagDepth int
)

func setupFlags() {
	flag.IntVar(&flagDepth, "depth", 1, "Depth of scraping")

	flag.Parse()
}
