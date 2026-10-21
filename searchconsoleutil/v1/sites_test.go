package searchconsoleutil

import (
	"context"
	"net/http"
	"testing"
)

func TestSites(t *testing.T) {
	f := newFake(t, map[string]func(http.ResponseWriter, []byte){
		"GET /webmasters/v3/sites": reply(`{"siteEntry":[
			{"siteUrl":"sc-domain:example.com","permissionLevel":"siteOwner"},
			{"siteUrl":"https://www.example.com/","permissionLevel":"siteRestrictedUser"},
			{"siteUrl":"https://other.test/","permissionLevel":"siteUnverifiedUser"}]}`),
	})
	got, err := f.client(t).Sites(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 {
		t.Fatalf("got %d sites", len(got))
	}
	if !got[0].IsDomainProperty() || got[1].IsDomainProperty() {
		t.Error("domain property not recognized")
	}
	if !got[1].CanRead() || got[2].CanRead() {
		t.Error("unverified user should not be able to read")
	}
}

func TestSitesEmpty(t *testing.T) {
	f := newFake(t, map[string]func(http.ResponseWriter, []byte){"GET /webmasters/v3/sites": reply(`{}`)})
	got, err := f.client(t).Sites(context.Background())
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("got %v, %v; want an empty, non-nil list", got, err)
	}
}

func TestSiteHostAndMatch(t *testing.T) {
	cases := []struct {
		site   string
		host   string
		domain string
		match  bool
	}{
		{"sc-domain:example.com", "example.com", "example.com", true},
		{"sc-domain:Example.COM", "example.com", "example.com", true},
		{"https://www.example.com/", "www.example.com", "example.com", true},
		{"https://www.example.com/", "www.example.com", "WWW.example.com.", true},
		{"sc-domain:example.com", "example.com", "blog.example.com", true},
		{"https://blog.example.com/", "blog.example.com", "example.com", true},
		{"sc-domain:example.com", "example.com", "notexample.com", false},
		{"https://notexample.com/", "notexample.com", "example.com", false},
		{"https://example.org/", "example.org", "example.com", false},
		{"https://example.com/", "example.com", "", false},
		{"://bad", "", "example.com", false},
	}
	for _, tc := range cases {
		s := Site{URL: tc.site}
		if got := s.Host(); got != tc.host {
			t.Errorf("Host(%q) = %q, want %q", tc.site, got, tc.host)
		}
		if got := s.MatchesDomain(tc.domain); got != tc.match {
			t.Errorf("MatchesDomain(%q, %q) = %v, want %v", tc.site, tc.domain, got, tc.match)
		}
	}
}
