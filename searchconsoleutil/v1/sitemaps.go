package searchconsoleutil

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Sitemap is a sitemap submitted for a property.
type Sitemap struct {
	// Path is the sitemap's URL.
	Path string
	// Type is the kind of sitemap, such as "sitemap" or "atomFeed".
	Type string
	// LastSubmitted and LastDownloaded are zero when Google gives none.
	LastSubmitted  time.Time
	LastDownloaded time.Time
	// IsPending is true while Google has not processed the sitemap.
	IsPending bool
	// IsIndex is true for a sitemap of sitemaps.
	IsIndex  bool
	Errors   int64
	Warnings int64
	Contents []SitemapContent
}

// SitemapContent counts a kind of URL in a sitemap.
type SitemapContent struct {
	Type      string
	Submitted int64
	Indexed   int64
}

// Sitemaps lists the sitemaps submitted for the property siteURL.
func (c *Client) Sitemaps(ctx context.Context, siteURL string) ([]Sitemap, error) {
	if siteURL == "" {
		return nil, errors.New("searchconsole: site URL cannot be empty")
	}
	resp, err := c.svc.Sitemaps.List(siteURL).Context(ctx).Do()
	if err != nil {
		return nil, wrap(err)
	}
	out := make([]Sitemap, 0, len(resp.Sitemap))
	for _, m := range resp.Sitemap {
		s := Sitemap{
			Path: m.Path, Type: m.Type, IsPending: m.IsPending, IsIndex: m.IsSitemapsIndex,
			Errors: m.Errors, Warnings: m.Warnings,
		}
		if s.LastSubmitted, err = parseTime(m.LastSubmitted); err != nil {
			return nil, fmt.Errorf("searchconsole: read sitemap %q submit time: %w", m.Path, err)
		}
		if s.LastDownloaded, err = parseTime(m.LastDownloaded); err != nil {
			return nil, fmt.Errorf("searchconsole: read sitemap %q download time: %w", m.Path, err)
		}
		for _, ct := range m.Contents {
			s.Contents = append(s.Contents, SitemapContent{Type: ct.Type, Submitted: ct.Submitted, Indexed: ct.Indexed})
		}
		out = append(out, s)
	}
	return out, nil
}

// parseTime reads an RFC 3339 time, or returns the zero time for "".
func parseTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, s)
}
