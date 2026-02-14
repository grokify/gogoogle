package slides

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/api/option"
	googleslides "google.golang.org/api/slides/v1"

	"github.com/grokify/gogoogle/cmd/gogoogle/internal/config"
	slidesutil "github.com/grokify/gogoogle/slidesutil/v1"
)

var (
	presentationID string
	includeNotes   bool
	prettyPrint    bool
)

var contentCmd = &cobra.Command{
	Use:   "content",
	Short: "Extract content from a Google Slides presentation",
	Long: `Extracts text and image URLs from a Google Slides presentation.

The output is JSON containing the presentation title, slide count, and for each slide:
- Title (if present)
- Text content from all elements
- Image URLs
- Speaker notes (if --notes flag is set)

Example:
  gogoogle slides content --presentation=1abc123xyz --notes --pretty`,
	RunE: runContent,
}

func init() {
	contentCmd.Flags().StringVarP(&presentationID, "presentation", "p", "",
		"Google Slides presentation ID (required)")
	contentCmd.Flags().BoolVarP(&includeNotes, "notes", "n", false,
		"Include speaker notes")
	contentCmd.Flags().BoolVar(&prettyPrint, "pretty", true,
		"Pretty print JSON output")

	_ = contentCmd.MarkFlagRequired("presentation")
}

func runContent(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	scopes := []string{
		googleslides.PresentationsReadonlyScope,
		googleslides.DriveReadonlyScope,
	}

	httpClient, err := config.NewHTTPClient(ctx, scopes)
	if err != nil {
		return fmt.Errorf("failed to create authenticated client: %w", err)
	}

	svc, err := googleslides.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return fmt.Errorf("failed to create Slides service: %w", err)
	}

	pres, err := svc.Presentations.Get(presentationID).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to get presentation: %w", err)
	}

	content := slidesutil.ExtractPresentationContent(pres, includeNotes)

	var output []byte
	if prettyPrint {
		output, err = json.MarshalIndent(content, "", "  ")
	} else {
		output, err = json.Marshal(content)
	}
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Fprintln(os.Stdout, string(output))
	return nil
}
