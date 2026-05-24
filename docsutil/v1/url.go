package docsutil

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// ErrInvalidURL is returned when a URL cannot be parsed.
	ErrInvalidURL = errors.New("invalid Google Docs URL")

	// ErrEmptyInput is returned when input is empty.
	ErrEmptyInput = errors.New("input cannot be empty")

	// rxDocsURL matches Google Docs document URLs.
	// Supports formats:
	//   - https://docs.google.com/document/d/{id}
	//   - https://docs.google.com/document/d/{id}/edit
	//   - https://docs.google.com/document/d/{id}/edit?tab=t.0
	//   - https://docs.google.com/document/d/{id}/edit?tab=t.0#heading=h.xyz
	//   - https://docs.google.com/document/d/{id}/edit#heading=h.xyz
	// The document ID is captured in group 1, ignoring any path suffix, query string, or fragment.
	rxDocsURL = regexp.MustCompile(`(?i)^https?://docs\.google\.com/document/d/([a-zA-Z0-9_-]+)(?:/[^?#]*)?(?:\?[^#]*)?(?:#.*)?$`)
)

// ParseDocumentURL extracts the document ID from a Google Docs URL.
// Returns the document ID or an error if the URL is invalid.
func ParseDocumentURL(urlOrID string) (string, error) {
	urlOrID = strings.TrimSpace(urlOrID)
	if urlOrID == "" {
		return "", ErrEmptyInput
	}

	// Check if it's already just an ID (no slashes or protocol)
	if !strings.Contains(urlOrID, "/") && !strings.Contains(urlOrID, ":") {
		return urlOrID, nil
	}

	// Try to parse as URL
	matches := rxDocsURL.FindStringSubmatch(urlOrID)
	if len(matches) < 2 {
		return "", ErrInvalidURL
	}

	return matches[1], nil
}

// BuildDocumentURL constructs a Google Docs URL from a document ID.
func BuildDocumentURL(documentID string) string {
	return "https://docs.google.com/document/d/" + documentID + "/edit"
}

// IsDocumentURL checks if the given string is a Google Docs document URL.
func IsDocumentURL(s string) bool {
	return rxDocsURL.MatchString(s)
}

// NormalizeDocumentInput accepts either a document ID or URL and returns the document ID.
// This is useful for APIs that accept either format.
func NormalizeDocumentInput(urlOrID string) (string, error) {
	return ParseDocumentURL(urlOrID)
}
