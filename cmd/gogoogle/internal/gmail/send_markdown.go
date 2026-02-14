package gmail

import (
	"context"
	"fmt"
	"net/mail"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/grokify/gogoogle/cmd/gogoogle/internal/config"
	gmailutil "github.com/grokify/gogoogle/gmailutil/v1"
	"github.com/grokify/mogo/mime/multipartutil"
	"github.com/grokify/mogo/net/mailutil"
)

var (
	// send-markdown command flags
	sendFrom    string
	sendTo      []string
	sendCc      []string
	sendBcc     []string
	sendSubject string
	sendBody    string
)

var sendMarkdownCmd = &cobra.Command{
	Use:   "send-markdown",
	Short: "Send an email with markdown body",
	Long: `Send an email with a markdown-formatted body.

The body can be specified as inline text or as a file reference using @filename.md.
The markdown is converted to both plain text and HTML for maximum compatibility.

Example:
  gogoogle gmail send-markdown \
    --goauth-credentials-file=creds.json \
    --goauth-credentials-account=myaccount \
    --to="user@example.com" \
    --subject="Hello" \
    --body="# Greeting\n\nHello world"

  # Or from a file:
  gogoogle gmail send-markdown \
    --to="user@example.com" \
    --subject="Newsletter" \
    --body=@newsletter.md`,
	RunE: runSendMarkdown,
}

func init() {
	sendMarkdownCmd.Flags().StringVar(&sendFrom, "from", "me",
		"Sender email or \"me\" for authenticated user")
	sendMarkdownCmd.Flags().StringSliceVar(&sendTo, "to", nil,
		"Recipient email addresses (required)")
	sendMarkdownCmd.Flags().StringSliceVar(&sendCc, "cc", nil,
		"CC email addresses")
	sendMarkdownCmd.Flags().StringSliceVar(&sendBcc, "bcc", nil,
		"BCC email addresses")
	sendMarkdownCmd.Flags().StringVarP(&sendSubject, "subject", "s", "",
		"Email subject (required)")
	sendMarkdownCmd.Flags().StringVarP(&sendBody, "body", "b", "",
		"Markdown body text or @filename.md to read from file (required)")

	_ = sendMarkdownCmd.MarkFlagRequired("to")
	_ = sendMarkdownCmd.MarkFlagRequired("subject")
	_ = sendMarkdownCmd.MarkFlagRequired("body")
}

func runSendMarkdown(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	scopes := []string{gmailutil.GmailSendScope}

	httpClient, err := config.NewHTTPClient(ctx, scopes)
	if err != nil {
		return fmt.Errorf("failed to create authenticated client: %w", err)
	}

	svc, err := gmailutil.NewGmailService(ctx, httpClient)
	if err != nil {
		return fmt.Errorf("failed to create Gmail service: %w", err)
	}

	// Parse body - handle @filename syntax.
	bodyText := sendBody
	if strings.HasPrefix(sendBody, "@") {
		filename := strings.TrimPrefix(sendBody, "@")
		data, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read body file %q: %w", filename, err)
		}
		bodyText = string(data)
	}

	// Parse recipient addresses.
	toAddrs := parseAddressList(sendTo)
	ccAddrs := parseAddressList(sendCc)
	bccAddrs := parseAddressList(sendBcc)

	// Create simple HTML from markdown (basic conversion).
	htmlBody := markdownToHTML(bodyText)

	// Build message.
	msg := mailutil.MessageWriter{
		To:      toAddrs,
		Cc:      ccAddrs,
		Bcc:     bccAddrs,
		Subject: sendSubject,
		BodyPartsSet: multipartutil.NewPartsSetAlternative(
			[]byte(bodyText), []byte(htmlBody)),
	}

	// Send the message.
	result, err := svc.Send(ctx, sendFrom, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Email sent successfully (ID: %s)\n", result.Id)
	return nil
}

func parseAddressList(addrs []string) mailutil.Addresses {
	var result mailutil.Addresses
	for _, addr := range addrs {
		parsed, err := mail.ParseAddress(addr)
		if err != nil {
			// Try as bare email.
			parsed = &mail.Address{Address: addr}
		}
		result = append(result, *parsed)
	}
	return result
}

// markdownToHTML performs a basic markdown to HTML conversion.
// For full markdown support, consider using a proper markdown library.
func markdownToHTML(md string) string {
	lines := strings.Split(md, "\n")
	var html strings.Builder
	html.WriteString("<html><body>\n")

	inParagraph := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Handle headers.
		if strings.HasPrefix(trimmed, "# ") {
			if inParagraph {
				html.WriteString("</p>\n")
				inParagraph = false
			}
			html.WriteString("<h1>")
			html.WriteString(escapeHTML(strings.TrimPrefix(trimmed, "# ")))
			html.WriteString("</h1>\n")
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			if inParagraph {
				html.WriteString("</p>\n")
				inParagraph = false
			}
			html.WriteString("<h2>")
			html.WriteString(escapeHTML(strings.TrimPrefix(trimmed, "## ")))
			html.WriteString("</h2>\n")
			continue
		}
		if strings.HasPrefix(trimmed, "### ") {
			if inParagraph {
				html.WriteString("</p>\n")
				inParagraph = false
			}
			html.WriteString("<h3>")
			html.WriteString(escapeHTML(strings.TrimPrefix(trimmed, "### ")))
			html.WriteString("</h3>\n")
			continue
		}

		// Handle empty lines.
		if trimmed == "" {
			if inParagraph {
				html.WriteString("</p>\n")
				inParagraph = false
			}
			continue
		}

		// Regular text - wrap in paragraph.
		if !inParagraph {
			html.WriteString("<p>")
			inParagraph = true
		} else {
			html.WriteString("<br>")
		}
		html.WriteString(escapeHTML(trimmed))
	}

	if inParagraph {
		html.WriteString("</p>\n")
	}
	html.WriteString("</body></html>")
	return html.String()
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
