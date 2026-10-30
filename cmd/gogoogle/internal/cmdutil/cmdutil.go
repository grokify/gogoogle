// Package cmdutil provides helpers shared by the gogoogle CLI commands:
// usage errors, argument validation, and JSON output.
package cmdutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// UsageError marks an error caused by how the command was invoked rather than
// by the API or the input files. The CLI exits with status 2 for these.
type UsageError struct{ Err error }

func (e *UsageError) Error() string { return e.Err.Error() }
func (e *UsageError) Unwrap() error { return e.Err }

// Usagef returns a UsageError with a formatted message.
func Usagef(format string, args ...any) error {
	return &UsageError{Err: fmt.Errorf(format, args...)}
}

// IsUsage reports whether err is, or wraps, a UsageError.
func IsUsage(err error) bool {
	var u *UsageError
	return errors.As(err, &u)
}

// Args wraps a cobra argument validator so its failure is a UsageError.
func Args(v cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := v(cmd, args); err != nil {
			return &UsageError{Err: err}
		}
		return nil
	}
}

// FlagError is a cobra flag error function that returns a UsageError.
func FlagError(_ *cobra.Command, err error) error { return &UsageError{Err: err} }

// WriteJSON writes v as indented JSON followed by a newline.
func WriteJSON(w io.Writer, v any) error {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}
	if _, err := fmt.Fprintln(w, string(out)); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}
