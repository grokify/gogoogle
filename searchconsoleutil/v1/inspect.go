package searchconsoleutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	searchconsole "google.golang.org/api/searchconsole/v1"
)

// Verdicts Search Console gives an inspection.
const (
	VerdictPass    = "PASS"
	VerdictPartial = "PARTIAL"
	VerdictFail    = "FAIL"
	VerdictNeutral = "NEUTRAL"
)

// Inspection is what Google knows about one URL in its index.
type Inspection struct {
	// URL is the page that was inspected.
	URL string
	// Verdict is one of the Verdict constants, or empty when Google gave no
	// index status.
	Verdict string
	// CoverageState is Google's description, such as "Submitted and indexed".
	CoverageState  string
	IndexingState  string
	PageFetchState string
	RobotsTxtState string
	// CrawledAs is DESKTOP or MOBILE.
	CrawledAs string
	// LastCrawl is zero when the page was never crawled.
	LastCrawl       time.Time
	GoogleCanonical string
	UserCanonical   string
	// Sitemaps are the sitemaps that list the page.
	Sitemaps      []string
	ReferringURLs []string
	// MobileUsabilityVerdict is empty when Google gave none.
	MobileUsabilityVerdict string
	// Link opens the inspection in the Search Console interface.
	Link string
}

// Indexed reports whether Google has the page in its index.
func (i Inspection) Indexed() bool { return i.Verdict == VerdictPass }

// InspectURL asks Google how pageURL is indexed. pageURL must be inside the
// property siteURL. Inspections are limited per property per day, so cache
// the results.
func (c *Client) InspectURL(ctx context.Context, siteURL, pageURL string) (*Inspection, error) {
	if siteURL == "" || pageURL == "" {
		return nil, errors.New("searchconsole: inspection needs a site URL and a page URL")
	}
	resp, err := c.svc.UrlInspection.Index.Inspect(&searchconsole.InspectUrlIndexRequest{
		SiteUrl: siteURL, InspectionUrl: pageURL,
	}).Context(ctx).Do()
	if err != nil {
		return nil, wrap(err)
	}
	out := &Inspection{URL: pageURL}
	r := resp.InspectionResult
	if r == nil {
		return out, nil
	}
	out.Link = r.InspectionResultLink
	if m := r.MobileUsabilityResult; m != nil {
		out.MobileUsabilityVerdict = m.Verdict
	}
	if ix := r.IndexStatusResult; ix != nil {
		out.Verdict = ix.Verdict
		out.CoverageState = ix.CoverageState
		out.IndexingState = ix.IndexingState
		out.PageFetchState = ix.PageFetchState
		out.RobotsTxtState = ix.RobotsTxtState
		out.CrawledAs = ix.CrawledAs
		out.GoogleCanonical = ix.GoogleCanonical
		out.UserCanonical = ix.UserCanonical
		out.Sitemaps = ix.Sitemap
		out.ReferringURLs = ix.ReferringUrls
		if out.LastCrawl, err = parseTime(ix.LastCrawlTime); err != nil {
			return nil, fmt.Errorf("searchconsole: read last crawl time: %w", err)
		}
	}
	return out, nil
}
