package driveutil

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/api/googleapi"
)

var (
	// ErrNotFound is returned when the Drive API responds with HTTP 404.
	ErrNotFound = errors.New("driveutil: not found")
	// ErrPermissionDenied is returned when the Drive API responds with HTTP 403.
	ErrPermissionDenied = errors.New("driveutil: permission denied")
	// ErrUnsupportedConversion is returned when a source file extension cannot
	// be converted to the requested Google file type.
	ErrUnsupportedConversion = errors.New("driveutil: unsupported conversion")
)

// wrapErr wraps err with the operation and resource ID. HTTP 404 and 403
// responses additionally match ErrNotFound and ErrPermissionDenied; the
// underlying *googleapi.Error stays reachable through errors.As.
func wrapErr(op, resource string, err error) error {
	if err == nil {
		return nil
	}
	var gerr *googleapi.Error
	if errors.As(err, &gerr) {
		switch gerr.Code {
		case http.StatusNotFound:
			return fmt.Errorf("driveutil: %s %s: %w: %w", op, resource, ErrNotFound, err)
		case http.StatusForbidden:
			return fmt.Errorf("driveutil: %s %s: %w: %w", op, resource, ErrPermissionDenied, err)
		}
	}
	return fmt.Errorf("driveutil: %s %s: %w", op, resource, err)
}
