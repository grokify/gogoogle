package docsutil

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	DocSlugDocument    = "document"
	DocSlugSpreadsheet = "spreadsheets"

	DocsURLFormat = `https://docs.google.com/%s/d/%s/edit#gid=0`
)

var rxDocURL = regexp.MustCompile(`(?i)^https:\/\/docs\.google\.com\/([^\/]+)\/d/([^\/]+)/`)

func ParseDocsURL(u, slug string) (id string, doctype string, err error) {
	u = strings.TrimSpace(u)
	slug = strings.ToLower(strings.TrimSpace(slug))
	if u == "" {
		return "", "", errors.New("url must be provided")
	}
	m := rxDocURL.FindAllStringSubmatch(u, -1)
	if len(m) == 0 {
		return "", "", errors.New("url cannot be parsed")
	}
	id = m[0][2]
	doctype = m[0][1]
	if slug != "" && slug != doctype {
		err = errors.New("doctype mismatch")
	}
	return
}

func BuildDocsURL(docsSlug, spreadsheetID string) string {
	return fmt.Sprintf(DocsURLFormat, docsSlug, spreadsheetID)
}
