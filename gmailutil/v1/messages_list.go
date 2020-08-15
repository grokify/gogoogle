package gmailutil

import (
	"strings"
	"time"

	"github.com/grokify/gotilla/time/timeutil"
)

const (
	GmailDateFormat  = "2006/01/02"
	ReferenceURL     = "https://developers.google.com/gmail/api/v1/reference"
	TutorialURLGO    = "https://developers.google.com/gmail/api/quickstart/go"
	ListApiReference = "https://developers.google.com/gmail/api/v1/reference/users/messages/list"
	ListApiExample   = "https://stackoverflow.com/questions/43057478/google-api-go-client-listing-messages-w-label-and-fetching-header-fields"
	FilteringExample = "https://developers.google.com/gmail/api/guides/filtering"
	FilterRules      = "https://support.google.com/mail/answer/7190"
)

// ?q=in:sent after:1388552400 before:1391230800
// ?q=in:sent after:2014/01/01 before:2014/02/01 ?q=in:sent after:2014/01/01 before:2014/02/01
// "from:someuser@example.com rfc822msgid:<somemsgid@example.com> is:unread"

/*

Warning: All dates used in the search query are interpretted as midnight on that date in the PST timezone. To specify accurate dates for other timezones pass the value in seconds instead:

*/

type MessageListQueryOpts struct {
	In          string
	From        string
	RFC822msgid string
	After       time.Time
	Before      time.Time
	Interval    timeutil.Interval
}

func GenerateMessageListQueryString(opts MessageListQueryOpts) string {
	parts := []string{}
	opts.From = strings.TrimSpace(opts.From)
	if len(opts.From) > 0 {
		parts = append(parts, "from:"+opts.From)
	}
	opts.In = strings.TrimSpace(opts.In)
	if len(opts.In) > 0 {
		parts = append(parts, "in:"+opts.In)
	}
	opts.RFC822msgid = strings.TrimSpace(opts.RFC822msgid)
	if len(opts.RFC822msgid) > 0 {
		parts = append(parts, "rfc822msgid:"+opts.RFC822msgid)
	}
	if !timeutil.TimeIsZeroAny(opts.After) {
		parts = append(parts, "after:"+opts.After.Format(GmailDateFormat))
	}
	if !timeutil.TimeIsZeroAny(opts.Before) {
		parts = append(parts, "before:"+opts.Before.Format(GmailDateFormat))
	}
	return strings.Join(parts, " ")
}
