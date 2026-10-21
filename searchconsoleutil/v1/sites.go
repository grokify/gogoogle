package searchconsoleutil

import (
	"context"
	"net/url"
	"strings"
)

// domainPrefix starts the identifier of a domain property.
const domainPrefix = "sc-domain:"

// Site is a Search Console property.
type Site struct {
	// URL identifies the property: "sc-domain:example.com" for a domain
	// property, or a URL such as "https://www.example.com/" for a
	// URL-prefix property. It is the value other methods take.
	URL string
	// PermissionLevel is siteOwner, siteFullUser, siteRestrictedUser or
	// siteUnverifiedUser.
	PermissionLevel string
}

// IsDomainProperty reports whether the property covers a whole domain,
// every subdomain and both protocols.
func (s Site) IsDomainProperty() bool { return strings.HasPrefix(s.URL, domainPrefix) }

// Host returns the lower-case host the property is for, or "" if it cannot
// be determined.
func (s Site) Host() string {
	if s.IsDomainProperty() {
		return strings.ToLower(strings.TrimPrefix(s.URL, domainPrefix))
	}
	u, err := url.Parse(s.URL)
	if err != nil {
		return ""
	}
	return strings.ToLower(u.Hostname())
}

// CanRead reports whether the credentials may query the property. An
// unverified user cannot.
func (s Site) CanRead() bool {
	return s.PermissionLevel != "" && s.PermissionLevel != "siteUnverifiedUser"
}

// MatchesDomain reports whether the property is about domain: its host is
// the domain, a subdomain of it, or a parent domain that covers it. domain is
// compared case-insensitively, and a leading "www." is not treated as
// special.
func (s Site) MatchesDomain(domain string) bool {
	domain = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(domain), "."))
	host := s.Host()
	if domain == "" || host == "" {
		return false
	}
	return host == domain ||
		strings.HasSuffix(host, "."+domain) ||
		strings.HasSuffix(domain, "."+host)
}

// Sites lists the properties the credentials can see.
func (c *Client) Sites(ctx context.Context) ([]Site, error) {
	resp, err := c.svc.Sites.List().Context(ctx).Do()
	if err != nil {
		return nil, wrap(err)
	}
	out := make([]Site, 0, len(resp.SiteEntry))
	for _, e := range resp.SiteEntry {
		out = append(out, Site{URL: e.SiteUrl, PermissionLevel: e.PermissionLevel})
	}
	return out, nil
}
