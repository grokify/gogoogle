package driveutil

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

// ConvertTo selects a Google file type to convert an uploaded file into.
type ConvertTo string

// Conversion targets. The empty value uploads the file unchanged.
const (
	ConvertNone   ConvertTo = ""
	ConvertDoc    ConvertTo = "doc"
	ConvertSheet  ConvertTo = "sheet"
	ConvertSlides ConvertTo = "slides"
)

// Google-native MIME types used as conversion targets.
const (
	MimeGoogleDoc    = "application/vnd.google-apps.document"
	MimeGoogleSheet  = "application/vnd.google-apps.spreadsheet"
	MimeGoogleSlides = "application/vnd.google-apps.presentation"
)

// conversions maps each target to its Google MIME type and the allowed
// source extensions with their media content types.
var conversions = map[ConvertTo]struct {
	mime string
	exts map[string]string
}{
	ConvertDoc: {MimeGoogleDoc, map[string]string{
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".odt":  "application/vnd.oasis.opendocument.text",
		".txt":  "text/plain",
		".html": "text/html",
	}},
	ConvertSheet: {MimeGoogleSheet, map[string]string{
		".csv":  "text/csv",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ods":  "application/vnd.oasis.opendocument.spreadsheet",
	}},
	ConvertSlides: {MimeGoogleSlides, map[string]string{
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".odp":  "application/vnd.oasis.opendocument.presentation",
	}},
}

// UploadRequest describes a file to upload.
type UploadRequest struct {
	Path      string    `json:"path"`      // local file; media type inferred from extension
	Name      string    `json:"name"`      // defaults to base name without extension when converting
	FolderID  string    `json:"folder_id"` // optional parent
	ConvertTo ConvertTo `json:"convert_to"`

	// Content, when set, is sent instead of opening Path. Path is then only a
	// name: its extension chooses the media type and the conversion check, and
	// its base name is the default Drive name. Use it to send a file the caller
	// has already opened and checked, so the path is not opened a second time.
	// Content is read once, to its end, and not closed.
	Content io.Reader `json:"-"`
}

// mediaPlan resolves the Drive metadata MIME type (empty when not
// converting) and the media content type for a path, validating the
// extension against the conversion target. It performs no I/O.
func mediaPlan(path string, convertTo ConvertTo) (metaMime, mediaType string, err error) {
	ext := strings.ToLower(filepath.Ext(path))
	if convertTo == ConvertNone {
		if t := mime.TypeByExtension(ext); t != "" {
			return "", t, nil
		}
		return "", "application/octet-stream", nil
	}
	c, ok := conversions[convertTo]
	if !ok {
		return "", "", fmt.Errorf("%w: unknown target %q", ErrUnsupportedConversion, convertTo)
	}
	mt, ok := c.exts[ext]
	if !ok {
		return "", "", fmt.Errorf("%w: extension %q cannot be converted to %q", ErrUnsupportedConversion, ext, convertTo)
	}
	return c.mime, mt, nil
}

// Upload creates a new Drive file from a local file, optionally converting it
// to a Google Docs, Sheets, or Slides file.
func (s *Service) Upload(ctx context.Context, r UploadRequest) (*File, error) {
	metaMime, mediaType, err := mediaPlan(r.Path, r.ConvertTo)
	if err != nil {
		return nil, fmt.Errorf("driveutil: upload %s: %w", r.Path, err)
	}
	name := r.Name
	if name == "" {
		name = filepath.Base(r.Path)
		if r.ConvertTo != ConvertNone {
			name = strings.TrimSuffix(name, filepath.Ext(name))
		}
	}
	meta := &drive.File{Name: name, MimeType: metaMime}
	if r.FolderID != "" {
		meta.Parents = []string{r.FolderID}
	}
	content := r.Content
	if content == nil {
		f, err := os.Open(r.Path)
		if err != nil {
			return nil, fmt.Errorf("driveutil: upload %s: %w", r.Path, err)
		}
		defer f.Close()
		content = f
	}
	out, err := s.DriveService.Files.Create(meta).
		Media(content, googleapi.ContentType(mediaType)).
		SupportsAllDrives(true).
		Fields(fileFields).
		Context(ctx).Do()
	if err != nil {
		return nil, wrapErr("upload", r.Path, err)
	}
	return fromDriveFile(out)
}

// Update replaces the content of an existing Drive file with a local file,
// keeping the file ID. A non-empty convertTo sets the Google MIME type.
func (s *Service) Update(ctx context.Context, fileID, path string, convertTo ConvertTo) (*File, error) {
	if _, _, err := mediaPlan(path, convertTo); err != nil {
		return nil, fmt.Errorf("driveutil: update %s: %w", fileID, err)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("driveutil: update %s: %w", fileID, err)
	}
	defer f.Close()
	return s.UpdateContent(ctx, fileID, path, f, convertTo)
}

// UpdateContent replaces the content of an existing Drive file with what
// content yields, keeping the file ID. name is only a name: its extension
// chooses the media type and the conversion check. content is read once, to its
// end, and not closed. Use it to send a file the caller has already opened and
// checked, as UploadRequest.Content does for Upload.
func (s *Service) UpdateContent(ctx context.Context, fileID, name string, content io.Reader, convertTo ConvertTo) (*File, error) {
	if content == nil {
		return nil, fmt.Errorf("driveutil: update %s: %w", fileID, errors.New("content is nil"))
	}
	metaMime, mediaType, err := mediaPlan(name, convertTo)
	if err != nil {
		return nil, fmt.Errorf("driveutil: update %s: %w", fileID, err)
	}
	out, err := s.DriveService.Files.Update(fileID, &drive.File{MimeType: metaMime}).
		Media(content, googleapi.ContentType(mediaType)).
		SupportsAllDrives(true).
		Fields(fileFields).
		Context(ctx).Do()
	if err != nil {
		return nil, wrapErr("update", fileID, err)
	}
	return fromDriveFile(out)
}
