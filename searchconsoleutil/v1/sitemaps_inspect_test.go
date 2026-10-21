package searchconsoleutil

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestSitemaps(t *testing.T) {
	f := newFake(t, map[string]func(http.ResponseWriter, []byte){
		"GET /webmasters/v3/sites/sc-domain:example.com/sitemaps": reply(`{"sitemap":[
			{"path":"https://example.com/sitemap.xml","type":"sitemap","lastSubmitted":"2026-10-01T10:00:00.000Z",
			 "lastDownloaded":"2026-10-18T04:30:00Z","isPending":false,"isSitemapsIndex":true,
			 "errors":"1","warnings":"2",
			 "contents":[{"type":"web","submitted":"120","indexed":"100"}]},
			{"path":"https://example.com/new.xml","isPending":true}]}`),
	})
	got, err := f.client(t).Sitemaps(context.Background(), "sc-domain:example.com")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d sitemaps", len(got))
	}
	a := got[0]
	if a.Path != "https://example.com/sitemap.xml" || !a.IsIndex || a.Errors != 1 || a.Warnings != 2 {
		t.Errorf("sitemap = %+v", a)
	}
	if a.LastSubmitted.IsZero() || a.LastDownloaded.Day() != 18 {
		t.Errorf("times = %v, %v", a.LastSubmitted, a.LastDownloaded)
	}
	if len(a.Contents) != 1 || a.Contents[0].Submitted != 120 || a.Contents[0].Indexed != 100 {
		t.Errorf("contents = %+v", a.Contents)
	}
	if b := got[1]; !b.IsPending || !b.LastSubmitted.IsZero() || !b.LastDownloaded.IsZero() {
		t.Errorf("pending sitemap = %+v", b)
	}
}

func TestSitemapsRejectsABadTime(t *testing.T) {
	f := newFake(t, map[string]func(http.ResponseWriter, []byte){
		"GET /webmasters/v3/sites/sc-domain:example.com/sitemaps": reply(`{"sitemap":[{"path":"x","lastSubmitted":"yesterday"}]}`),
	})
	if _, err := f.client(t).Sitemaps(context.Background(), "sc-domain:example.com"); err == nil {
		t.Fatal("want an error")
	}
	if _, err := f.client(t).Sitemaps(context.Background(), ""); err == nil {
		t.Fatal("empty site URL: want an error")
	}
}

const inspectPath = "POST /v1/urlInspection/index:inspect"

func TestInspectURL(t *testing.T) {
	var sent map[string]string
	f := newFake(t, map[string]func(http.ResponseWriter, []byte){
		inspectPath: func(w http.ResponseWriter, body []byte) {
			if err := json.Unmarshal(body, &sent); err != nil {
				t.Error(err)
			}
			if _, err := w.Write([]byte(`{"inspectionResult":{
				"inspectionResultLink":"https://search.google.com/search-console/inspect?x=1",
				"indexStatusResult":{"verdict":"PASS","coverageState":"Submitted and indexed",
					"indexingState":"INDEXING_ALLOWED","pageFetchState":"SUCCESSFUL","robotsTxtState":"ALLOWED",
					"crawledAs":"MOBILE","lastCrawlTime":"2026-10-17T08:00:00Z",
					"googleCanonical":"https://example.com/","userCanonical":"https://example.com/",
					"sitemap":["https://example.com/sitemap.xml"],"referringUrls":["https://example.com/a"]},
				"mobileUsabilityResult":{"verdict":"PASS"}}}`)); err != nil {
				t.Error(err)
			}
		},
	})
	got, err := f.client(t).InspectURL(context.Background(), "sc-domain:example.com", "https://example.com/")
	if err != nil {
		t.Fatal(err)
	}
	if sent["siteUrl"] != "sc-domain:example.com" || sent["inspectionUrl"] != "https://example.com/" {
		t.Errorf("request = %v", sent)
	}
	if !got.Indexed() || got.CoverageState != "Submitted and indexed" || got.CrawledAs != "MOBILE" {
		t.Errorf("inspection = %+v", got)
	}
	if got.LastCrawl.Day() != 17 || got.GoogleCanonical != "https://example.com/" || len(got.Sitemaps) != 1 {
		t.Errorf("inspection = %+v", got)
	}
	if got.MobileUsabilityVerdict != VerdictPass || got.Link == "" || got.URL != "https://example.com/" {
		t.Errorf("inspection = %+v", got)
	}
}

func TestInspectURLWithoutAnIndexStatus(t *testing.T) {
	for name, body := range map[string]string{
		"no result": `{}`,
		"no status": `{"inspectionResult":{"inspectionResultLink":"l"}}`,
	} {
		f := newFake(t, map[string]func(http.ResponseWriter, []byte){inspectPath: reply(body)})
		got, err := f.client(t).InspectURL(context.Background(), "sc-domain:example.com", "https://example.com/x")
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if got.Indexed() || got.Verdict != "" || !got.LastCrawl.IsZero() {
			t.Errorf("%s: inspection = %+v", name, got)
		}
	}
}

func TestInspectURLValidationAndQuota(t *testing.T) {
	f := newFake(t, map[string]func(http.ResponseWriter, []byte){inspectPath: fail(http.StatusTooManyRequests, "rateLimitExceeded")})
	c := f.client(t)
	if _, err := c.InspectURL(context.Background(), "", "https://example.com/"); err == nil {
		t.Error("empty site: want an error")
	}
	if _, err := c.InspectURL(context.Background(), "sc-domain:example.com", ""); err == nil {
		t.Error("empty page: want an error")
	}
	if len(f.paths) != 0 {
		t.Errorf("invalid input reached the API: %v", f.paths)
	}
	_, err := c.InspectURL(context.Background(), "sc-domain:example.com", "https://example.com/")
	if !errors.Is(err, ErrQuotaExceeded) {
		t.Fatalf("err = %v, want ErrQuotaExceeded", err)
	}
}
