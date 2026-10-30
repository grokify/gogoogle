// Package sheets provides the sheets group command for the gogoogle CLI.
package sheets

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/api/option"
	gsheets "google.golang.org/api/sheets/v4"

	"github.com/grokify/gogoogle/cmd/gogoogle/internal/cmdutil"
	"github.com/grokify/gogoogle/cmd/gogoogle/internal/config"
	sheetsutil "github.com/grokify/gogoogle/sheetsutil/v4"
)

// Cmd is the sheets group command.
var Cmd = &cobra.Command{
	Use:   "sheets",
	Short: "Google Sheets utilities",
	Long:  `Commands for reading Google Sheets.`,
}

// open creates the Sheets service. Tests replace it.
var open = func(ctx context.Context) (*gsheets.Service, error) {
	c, err := config.NewHTTPClient(ctx, []string{gsheets.SpreadsheetsReadonlyScope})
	if err != nil {
		return nil, fmt.Errorf("create authenticated client: %w", err)
	}
	return gsheets.NewService(ctx, option.WithHTTPClient(c))
}

func init() {
	Cmd.AddCommand(recordsCmd())
}

func recordsCmd() *cobra.Command {
	var sheet string
	cmd := &cobra.Command{
		Use:   "records <spreadsheetID|URL>",
		Short: "Read a sheet as records keyed by its header row",
		Long: `Reads one sheet and prints a JSON array with one object per row, keyed by the
header in row 1. Values are the formatted text shown in the sheet. Rows with
no values are skipped; a repeated header is an error.

Needs the spreadsheets.readonly scope.

Example:
  gogoogle sheets records 1AbC --sheet "Round 1"`,
		Args: cmdutil.Args(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := sheetsutil.ParseSpreadsheetURL(args[0])
			if err != nil {
				return cmdutil.Usagef("spreadsheet: %w", err)
			}
			svc, err := open(cmd.Context())
			if err != nil {
				return err
			}
			records, err := sheetsutil.ReadRecords(cmd.Context(), svc, id, sheet)
			if err != nil {
				return err
			}
			return cmdutil.WriteJSON(cmd.OutOrStdout(), records)
		},
	}
	cmd.Flags().StringVar(&sheet, "sheet", "", "Sheet name (default: the first sheet)")
	return cmd
}
