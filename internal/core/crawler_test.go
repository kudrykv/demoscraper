package core_test

import (
	"context"
	"demoscraper/internal/adapters/webpager"
	"demoscraper/internal/clients/xresty"
	"demoscraper/internal/core"
	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
	"net/http"
	"testing"
	"time"
)

func TestCrawler_Crawl(t *testing.T) {
	t.Parallel()

	t.Run("successful run", func(t *testing.T) {
		t.Parallel()

		vcr := setupVCR("fixtures/crawler/successful_run")
		defer func() { require.NoError(t, vcr.Stop()) }()

		ctx := context.Background()

		rawHTTPClient := &http.Client{Transport: vcr}
		httpClient := xresty.New(rawHTTPClient)
		webPager := webpager.New(httpClient)
		crawler := core.NewCrawler(webPager)

		crawlEntries, err := crawler.Crawl(ctx, core.CrawlParameters{StartURL: "https://github.com", DepthLimit: 2})
		require.NoError(t, err)

		for {
			timer := time.NewTimer(5 * time.Second)

			select {
			case crawlEntry, ok := <-crawlEntries:
				if !ok {
					return
				}

				t.Logf("CrawlEntry: %+v", crawlEntry)
				timer.Stop()
			case <-timer.C:
				require.FailNow(t, "timeout while waiting for crawl entry")
			}
		}
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
