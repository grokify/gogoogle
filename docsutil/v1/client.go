// Package docsutil provides utilities for working with the Google Docs API.
package docsutil

import (
	"context"
	"net/http"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

// Scopes returns the OAuth2 scopes for Google Docs API access.
func Scopes() []string {
	return []string{
		docs.DocumentsReadonlyScope,
		docs.DriveReadonlyScope,
	}
}

// ScopesReadWrite returns the OAuth2 scopes for read-write Google Docs API access.
func ScopesReadWrite() []string {
	return []string{
		docs.DocumentsScope,
		docs.DriveScope,
	}
}

// Service wraps the Google Docs API service.
type Service struct {
	httpClient       *http.Client
	DocsService      *docs.Service
	DocumentsService *docs.DocumentsService
}

// NewService creates a new Docs service from an authenticated HTTP client.
func NewService(ctx context.Context, httpClient *http.Client) (*Service, error) {
	svc, err := docs.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	return &Service{
		httpClient:       httpClient,
		DocsService:      svc,
		DocumentsService: docs.NewDocumentsService(svc),
	}, nil
}

// GetDocument retrieves a document by ID.
func (s *Service) GetDocument(ctx context.Context, documentID string) (*docs.Document, error) {
	return s.DocumentsService.Get(documentID).Context(ctx).Do()
}

// GetDocumentWithSuggestions retrieves a document with the specified suggestions view mode.
// Mode can be: DEFAULT_FOR_CURRENT_ACCESS, SUGGESTIONS_INLINE, PREVIEW_SUGGESTIONS_ACCEPTED,
// PREVIEW_WITHOUT_SUGGESTIONS
func (s *Service) GetDocumentWithSuggestions(ctx context.Context, documentID, suggestionsViewMode string) (*docs.Document, error) {
	return s.DocumentsService.Get(documentID).
		SuggestionsViewMode(suggestionsViewMode).
		Context(ctx).Do()
}
