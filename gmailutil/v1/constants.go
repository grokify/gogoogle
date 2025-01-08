package gmailutil

import gmail "google.golang.org/api/gmail/v1"

// https://pkg.go.dev/google.golang.org/api/gmail/v1#pkg-constants

const (
	MailGoogleComScope = gmail.MailGoogleComScope // "https://mail.google.com/"
	GmailReadonlyScope = gmail.GmailReadonlyScope // "https://www.googleapis.com/auth/gmail.readonly"
	GmailSendScope     = gmail.GmailSendScope     // "https://www.googleapis.com/auth/gmail.send"

	UserIDMe = "me"
)
