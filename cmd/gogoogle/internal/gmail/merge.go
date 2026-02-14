package gmail

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/authutil"
	"github.com/grokify/gogoogle/gmailutil/v1/mailmerge"
)

var (
	// merge command flags
	mergeSheetID         string
	mergeSheetIndex      uint
	mergeSheetHeaderRows uint32
	mergeSubjectTemplate string
	mergeHTMLTemplate    string
	mergeTextTemplate    string
	mergeInlineFiles     []string
	mergeAttachmentFiles []string
	mergeGoauthFile      string
	mergeGoauthAccount   string
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Send templated emails via mail merge",
	Long: `Send templated emails using data from a Google Sheet.

The Google Sheet should contain columns for recipients (TO, CC, BCC) and any
template variables used in the subject and body templates.

Template files use Mustache syntax. The template variables are populated from
the column headers in the Google Sheet.

Example:
  gogoogle gmail merge \
    --goauth-credentials-file=creds.json \
    --goauth-credentials-account=myaccount \
    --sheet-id=1abc123xyz \
    --subject-template=subject.mustache \
    --html-template=body.mustache`,
	RunE: runMerge,
}

func init() {
	mergeCmd.Flags().StringVarP(&mergeGoauthFile, "goauth-credentials-file", "c", "",
		"Path to goauth CredentialsSet JSON file (env: GOAUTH_CREDENTIALS_FILE)")
	mergeCmd.Flags().StringVarP(&mergeGoauthAccount, "goauth-credentials-account", "k", "",
		"Account key within goauth CredentialsSet file (env: GOAUTH_CREDENTIALS_ACCOUNT)")
	mergeCmd.Flags().StringVarP(&mergeSheetID, "sheet-id", "s", "",
		"Google Sheet ID with recipients (required)")
	mergeCmd.Flags().UintVarP(&mergeSheetIndex, "sheet-index", "x", 0,
		"Sheet index within the spreadsheet")
	mergeCmd.Flags().Uint32VarP(&mergeSheetHeaderRows, "sheet-header-rows", "r", 1,
		"Number of header rows in the sheet")
	mergeCmd.Flags().StringVarP(&mergeSubjectTemplate, "subject-template", "j", "",
		"Subject template file (required)")
	mergeCmd.Flags().StringVar(&mergeHTMLTemplate, "html-template", "",
		"HTML body template file")
	mergeCmd.Flags().StringVarP(&mergeTextTemplate, "text-template", "t", "",
		"Text body template file")
	mergeCmd.Flags().StringSliceVarP(&mergeInlineFiles, "inline", "i", nil,
		"Inline attachment files")
	mergeCmd.Flags().StringSliceVarP(&mergeAttachmentFiles, "attachment", "a", nil,
		"Attachment files")

	_ = mergeCmd.MarkFlagRequired("sheet-id")
	_ = mergeCmd.MarkFlagRequired("subject-template")
}

func runMerge(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Apply environment variable defaults.
	if mergeGoauthFile == "" {
		mergeGoauthFile = os.Getenv("GOAUTH_CREDENTIALS_FILE")
	}
	if mergeGoauthAccount == "" {
		mergeGoauthAccount = os.Getenv("GOAUTH_CREDENTIALS_ACCOUNT")
	}

	if mergeGoauthFile == "" || mergeGoauthAccount == "" {
		return fmt.Errorf("goauth credentials required: use --goauth-credentials-file and --goauth-credentials-account")
	}

	var googleClient *http.Client
	creds, err := goauth.NewCredentialsFromSetFile(mergeGoauthFile, mergeGoauthAccount, true)
	if err != nil {
		return fmt.Errorf("failed to load credentials: %w", err)
	}
	tok, err := creds.NewOrExistingValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	googleClient = authutil.NewClientTokenOAuth2(tok)

	opts := mailmerge.MailMergeOpts{
		GoauthCredsFile:                 mergeGoauthFile,
		GoauthAccountKey:                mergeGoauthAccount,
		RecipientsGoogleSheetID:         mergeSheetID,
		RecipientsGoogleSheetIndex:      mergeSheetIndex,
		RecipientsGoogleSheetHeaderRows: mergeSheetHeaderRows,
		SubjectTemplateTextFilename:     mergeSubjectTemplate,
		BodyTemplateHTMLFilename:        mergeHTMLTemplate,
		BodyTemplateTextFilename:        mergeTextTemplate,
		InlineFilenames:                 mergeInlineFiles,
		AttachmentsFilenames:            mergeAttachmentFiles,
		GoogleClient:                    googleClient,
	}

	mm, err := mailmerge.NewMailMerge(ctx, &opts)
	if err != nil {
		return fmt.Errorf("failed to create mail merge: %w", err)
	}

	cnt, err := mm.Send(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to send mail merge: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Successfully sent %d email message(s)\n", cnt)
	return nil
}
