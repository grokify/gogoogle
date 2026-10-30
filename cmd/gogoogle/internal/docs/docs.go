// Package docs provides the docs group command for the gogoogle CLI.
package docs

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/grokify/gogoogle/cmd/gogoogle/internal/cmdutil"
	"github.com/grokify/gogoogle/cmd/gogoogle/internal/config"
	docsutil "github.com/grokify/gogoogle/docsutil/v1"
)

// Cmd is the docs group command.
var Cmd = &cobra.Command{
	Use:   "docs",
	Short: "Google Docs utilities",
	Long:  `Commands for reading Google Docs, including pending suggestions.`,
}

// open creates the Docs service. Tests replace it.
var open = func(ctx context.Context) (*docsutil.Service, error) {
	c, err := config.NewHTTPClient(ctx, docsutil.Scopes())
	if err != nil {
		return nil, fmt.Errorf("create authenticated client: %w", err)
	}
	return docsutil.NewService(ctx, c)
}

func init() {
	Cmd.AddCommand(suggestionsCmd())
}

func suggestionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "suggestions <documentID|URL>",
		Short: "Show a document's text split into original, inserted, and deleted segments",
		Long: `Fetches a document with its suggestions shown inline and prints each paragraph
and table cell as JSON, with the text classified as original, insertion, or
deletion, plus the original and as-if-accepted text of the block.

The Docs API does not report who made a suggestion, so the output has no
authors; comments (gogoogle drive comments) do. Needs a Docs read scope.

Example:
  gogoogle docs suggestions 1AbC`,
		Args: cmdutil.Args(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := docsutil.NormalizeDocumentInput(args[0])
			if err != nil {
				return cmdutil.Usagef("document: %w", err)
			}
			svc, err := open(cmd.Context())
			if err != nil {
				return err
			}
			doc, err := svc.GetDocumentWithSuggestions(cmd.Context(), id, "SUGGESTIONS_INLINE")
			if err != nil {
				return fmt.Errorf("get document %s: %w", id, err)
			}
			return cmdutil.WriteJSON(cmd.OutOrStdout(), docsutil.ExtractSuggestionContent(doc))
		},
	}
}
