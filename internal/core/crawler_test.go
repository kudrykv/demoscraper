package core_test

import (
	"context"
	"demoscraper/internal/adapters/inmemvisitor"
	"demoscraper/internal/adapters/webpager"
	"demoscraper/internal/clients/xresty"
	"demoscraper/internal/core"
	"demoscraper/internal/core/entities"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

func TestCrawler_Crawl(t *testing.T) {
	t.Parallel()

	uniqueURLsNumber := 2854

	t.Run("successful run", func(t *testing.T) {
		t.Parallel()

		vcr := setupVCR("fixtures/crawler/successful_run")
		defer func() { require.NoError(t, vcr.Stop()) }()

		crawler := setupCrawler(vcr)

		ctx := context.Background()

		crawlParameters := core.CrawlParameters{StartURL: "https://github.com", DepthLimit: 2, Parallelism: 8}
		crawlEntries, err := crawler.Crawl(ctx, crawlParameters)

		require.NoError(t, err)

		urls := drainCrawlEntries(t, crawlEntries)

		require.Equal(t, uniqueURLsNumber, len(urls))

		urlMap := make(map[string]struct{})
		for _, url := range urls {
			urlMap[url] = struct{}{}
		}

		require.Equal(t, len(urls), len(urlMap))
	})

	t.Run("preliminary cancel should stop crawling", func(t *testing.T) {
		t.Parallel()

		vcr := setupVCR("fixtures/crawler/successful_run")
		defer func() { require.NoError(t, vcr.Stop()) }()

		crawler := setupCrawler(vcr)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		crawlParameters := core.CrawlParameters{StartURL: "https://github.com", DepthLimit: 2, Parallelism: 8}
		crawlEntries, err := crawler.Crawl(ctx, crawlParameters)

		require.NoError(t, err)

		urls := drainCrawlEntries(t, crawlEntries)

		require.Equal(t, 0, len(urls))
	})

	t.Run("cancel after the crawl should stop crawling", func(t *testing.T) {
		t.Parallel()

		vcr := setupVCR("fixtures/crawler/successful_run")
		defer func() { require.NoError(t, vcr.Stop()) }()

		crawler := setupCrawler(vcr)

		ctx, cancel := context.WithCancel(context.Background())

		crawlParameters := core.CrawlParameters{StartURL: "https://github.com", DepthLimit: 2, Parallelism: 8}
		crawlEntries, err := crawler.Crawl(ctx, crawlParameters)

		require.NoError(t, err)

		cancel()

		urls := drainCrawlEntriesWithTimeout(t, crawlEntries, time.Second)
		require.Less(t, len(urls), uniqueURLsNumber)
	})
}

func setupVCR(cassetteName string) *recorder.Recorder {
	vcr, err := recorder.NewWithOptions(&recorder.Options{
		Mode:         recorder.ModeRecordOnce,
		CassetteName: cassetteName,
	})
	if err != nil {
		panic(err)
	}

	return vcr
}

func setupCrawler(roundTripper http.RoundTripper) *core.Crawler {
	rawHTTPClient := &http.Client{Transport: roundTripper}
	httpClient := xresty.New(rawHTTPClient)
	visitorMaker := inmemvisitor.New
	webPager := webpager.New(httpClient)

	return core.NewCrawler(webPager, visitorMaker)
}

func drainCrawlEntries(t *testing.T, crawlEntries <-chan entities.CrawlEntry) []string {
	t.Helper()

	var urls []string

	for {
		timer := time.NewTimer(5 * time.Second)
		finished := false

		select {
		case crawlEntry, ok := <-crawlEntries:
			if !ok {
				finished = true

				break
			}

			urls = append(urls, crawlEntry.URL())

			timer.Stop()
		case <-timer.C:
			require.FailNow(t, "timeout while waiting for crawl entry")
		}

		if finished {
			break
		}
	}

	return urls
}

func drainCrawlEntriesWithTimeout(
	t *testing.T, crawlEntries <-chan entities.CrawlEntry, timeout time.Duration,
) []string {
	t.Helper()

	var urls []string

	timeoutTimer := time.NewTimer(timeout)

	for {
		stop := false

		select {
		case crawlEntry, ok := <-crawlEntries:
			if !ok {
				timeoutTimer.Stop()

				stop = true

				break
			}

			urls = append(urls, crawlEntry.URL())
		case <-timeoutTimer.C:
			require.FailNow(t, "timeout while waiting for crawl entry")
		}

		if stop {
			break
		}
	}

	return urls
}
