// Example mail_merge demonstrates sending templated emails using the mailmerge package.
//
// This example shows how to:
//   - Use Mustache templates for email subject and body
//   - Include inline images via Content-ID references
//   - Read recipient data from a Google Sheet
//
// Usage:
//
//	go run main.go \
//	  --goauth-credentials-file=/path/to/credentials.json \
//	  --goauth-account-key=myaccount \
//	  --sheet-id=YOUR_GOOGLE_SHEET_ID
//
// The Google Sheet should have columns: TO, CC, BCC, FIRST_NAME, COMPANY_NAME, MONTH,
// HEADLINE, SUMMARY, ARTICLE_URL, UNSUBSCRIBE_URL
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/grokify/gogoogle/gmailutil/v1/mailmerge"
)

func main() {
	cnt, err := mailmerge.ExecMailMergeCLI(context.Background())
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	fmt.Printf("Successfully sent %d email message(s)\n", cnt)
}
