// Package driveutil provides utilities for working with the Google Drive API v3.
package driveutil

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const fileFields = "id,name,mimeType,webViewLink,modifiedTime,headRevisionId"

// Service wraps the Google Drive API service.
type Service struct {
	httpClient   *http.Client
	DriveService *drive.Service
}

// NewService creates a new Drive service from an authenticated HTTP client.
// Additional options (for example option.WithEndpoint) are applied after the
// HTTP client.
func NewService(ctx context.Context, c *http.Client, opts ...option.ClientOption) (*Service, error) {
	all := make([]option.ClientOption, 0, len(opts)+1)
	all = append(all, option.WithHTTPClient(c))
	all = append(all, opts...)
	svc, err := drive.NewService(ctx, all...)
	if err != nil {
		return nil, fmt.Errorf("driveutil: new service: %w", err)
	}
	return &Service{httpClient: c, DriveService: svc}, nil
}

// File is the subset of Drive file metadata used by this package.
type File struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	MimeType       string    `json:"mime_type"`
	WebViewLink    string    `json:"web_view_link"`
	HeadRevisionID string    `json:"head_revision_id"`
	ModifiedTime   time.Time `json:"modified_time"`
}

// GetFile retrieves file metadata by ID.
func (s *Service) GetFile(ctx context.Context, fileID string) (*File, error) {
	f, err := s.DriveService.Files.Get(fileID).
		Fields(fileFields).
		SupportsAllDrives(true).
		Context(ctx).Do()
	if err != nil {
		return nil, wrapErr("get file", fileID, err)
	}
	return fromDriveFile(f)
}

func fromDriveFile(f *drive.File) (*File, error) {
	out := &File{
		ID:             f.Id,
		Name:           f.Name,
		MimeType:       f.MimeType,
		WebViewLink:    f.WebViewLink,
		HeadRevisionID: f.HeadRevisionId,
	}
	t, err := parseTime(f.ModifiedTime)
	if err != nil {
		return nil, fmt.Errorf("driveutil: file %s: modifiedTime: %w", f.Id, err)
	}
	out.ModifiedTime = t
	return out, nil
}

func parseTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, s)
}
