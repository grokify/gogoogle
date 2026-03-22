# Sending Emails

## SendSimple

The `SendSimple` method provides the easiest way to send emails:

```go
import "github.com/grokify/gogoogle/gmailutil/v1"

service, _ := gmailutil.NewGmailService(ctx, httpClient)

_, err := service.SendSimple(ctx, "me", gmailutil.SendSimpleOpts{
    To:       "recipient@example.com",
    Subject:  "Meeting Tomorrow",
    BodyText: "Let's meet at 2pm.",
    BodyHTML: "<p>Let's meet at <b>2pm</b>.</p>",
    ReplyTo:  "replies@mycompany.com",
})
```

### SendSimpleOpts

| Field | Type | Description |
|-------|------|-------------|
| `To` | `string` | Recipient email address |
| `Subject` | `string` | Email subject line |
| `BodyText` | `string` | Plain text body (required) |
| `BodyHTML` | `string` | HTML body (optional, creates multipart) |
| `ReplyTo` | `string` | Reply-To header (optional) |

### Text Only vs HTML

```go
// Text only
opts := gmailutil.SendSimpleOpts{
    To:       "user@example.com",
    Subject:  "Plain Text Email",
    BodyText: "This is plain text.",
}

// Text + HTML (multipart/alternative)
opts := gmailutil.SendSimpleOpts{
    To:       "user@example.com",
    Subject:  "Rich Email",
    BodyText: "This is the plain text fallback.",
    BodyHTML: "<h1>Rich Email</h1><p>This is HTML content.</p>",
}
```

## Send with MessageWriter

For advanced use cases, use `mailutil.MessageWriter`:

```go
import (
    "github.com/grokify/mogo/net/mailutil"
    "github.com/grokify/gogoogle/gmailutil/v1"
)

msg := mailutil.MessageWriter{
    From:    mailutil.Address{Address: "sender@example.com", Name: "Sender Name"},
    To:      mailutil.Addresses{{Address: "recipient@example.com", Name: "Recipient"}},
    Cc:      mailutil.Addresses{{Address: "cc@example.com"}},
    Subject: "Hello",
    Header: map[string][]string{
        "X-Custom-Header": {"custom-value"},
    },
}

// Set body
msg.BodyPartsSet, _ = multipartutil.NewPartsSetMail(
    []byte("Plain text"),
    []byte("<p>HTML</p>"),
    nil, // attachments
)

_, err := service.Send(ctx, "me", msg)
```

### Multiple Recipients

```go
msg := mailutil.MessageWriter{
    To: mailutil.Addresses{
        {Address: "user1@example.com", Name: "User One"},
        {Address: "user2@example.com", Name: "User Two"},
    },
    Cc: mailutil.Addresses{
        {Address: "cc@example.com"},
    },
    Bcc: mailutil.Addresses{
        {Address: "bcc@example.com"},
    },
    Subject: "Group Email",
}
```

### Attachments

```go
import "github.com/grokify/mogo/mime/multipartutil"

// Read file
fileData, _ := os.ReadFile("report.pdf")

// Create parts set with attachment
partsSet, _ := multipartutil.NewPartsSetMail(
    []byte("Please see attached report."),
    []byte("<p>Please see attached report.</p>"),
    []multipartutil.Part{
        {
            ContentType: "application/pdf",
            Filename:    "report.pdf",
            Data:        fileData,
        },
    },
)

msg.BodyPartsSet = partsSet
```

## From Address

The `from` parameter specifies the sender:

```go
// Send as authenticated user
service.SendSimple(ctx, "me", opts)

// Send as specific address (must have permission)
service.SendSimple(ctx, "alias@example.com", opts)
```

Use `"me"` to send from the authenticated user's primary address.

## Error Handling

```go
result, err := service.SendSimple(ctx, "me", opts)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "Invalid To"):
        log.Println("Invalid recipient address")
    case strings.Contains(err.Error(), "Daily Limit"):
        log.Println("Gmail sending limit reached")
    default:
        log.Printf("Send error: %v", err)
    }
    return
}

log.Printf("Email sent, ID: %s", result.Id)
```

## Gmail Sending Limits

| Account Type | Daily Limit |
|--------------|-------------|
| Gmail (free) | 500 emails |
| Google Workspace | 2,000 emails |

For bulk sending, consider:

- Batching with delays
- Using mail merge for personalization
- Google Workspace for higher limits

## Next Steps

- [Mail Merge](mail-merge.md) - Template-based bulk sending
- [Messages](messages.md) - Read and manage emails
