package main

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/grokify/goauth"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/mime/multipartutil"
	"github.com/grokify/mogo/net/mailutil"

	"github.com/grokify/gogoogle/gmailutil/v1"
)

func main() {
	opts, err := goauth.ParseOptions()
	logutil.FatalErr(err)

	client, err := opts.NewClient(context.Background(), "mystate")
	logutil.FatalErr(err)

	svc, err := gmailutil.NewGmailService(context.Background(), client)
	logutil.FatalErr(err)

	emailAddrAlice := mail.Address{Address: "alice@example.com", Name: "Alice"}
	emailAddrBob := mail.Address{Address: "bob@example.com", Name: "Bob"}

	emBodyText := fmt.Sprintf(`Hi %s!`, emailAddrBob.Name)
	emBodyHTML := fmt.Sprintf(`<html><body>Hi <a href="mailto:%s">%s</a>!<body></html>`, emailAddrBob.Address, emailAddrBob.Name)

	msg, err := svc.Send(context.Background(), emailAddrAlice.Address, mailutil.MessageWriter{
		To:      mailutil.Addresses{emailAddrAlice},
		Subject: "At Mention API Test",
		BodyPartsSet: multipartutil.NewPartsSetAlternative(
			[]byte(emBodyText), []byte(emBodyHTML)),
	})
	logutil.FatalErr(err)

	fmt.Printf("SENT EMAIL ID (%s)\n", msg.Id)
}
