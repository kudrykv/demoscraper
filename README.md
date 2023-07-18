# Demo scraper

The app is in a form of a simple CLI.

## How to run

```shell
go run ./cmd/demoscraper -url https://github.com
```

If URL is not provided, the app will use `https://github.com` as a default.
Another available flags:
- `--depth` - how many levels of links to scrape (default: 2)
- `--output` - output file name (default: `output.tsv`)
- `--parallelism` - how many concurrent requests to make (default: number of CPU)

## Description

Crawler logic is straightforward: it starts with a given URL, fetches the page, parses it, and then
iteratively does the same for all the links found on the page. It also keeps track of the depth of
the current link, so it can stop when it reaches the maximum depth.

The app uses a simple in-memory cache to avoid fetching the same URL twice, and
a semaphore to limit the number of concurrent requests.

### Structure

The app is split into two large packages: `cmd` and `internal`.
The former contains the CLI app, the latter contains the app internals.

App internals are split into several packages:
- `core` — contains the main crawler logic
- `adapters` – contains the implementations of the interfaces defined in `core` package
- `clients` — contains the wrappers for external libraries' clients

`./internal/core/crawler.go` is the essence of the app.
It contains the main crawler logic.
It iterates through the depth levels.
On each level, it parallels page fetching and parsing.
With a simple in-memory cache it avoids fetching the same URL twice.
Once all the links of the given depth level are fetched and parsed, it moves to the next level.

### Tests

For the sake of simplicity, the app has only a few tests.
The main crawler logic is covered with unit tests.
It checks happy paths, as well as handling of cancellation signal.
