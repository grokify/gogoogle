# Gmail

The `gmailutil/v1` package provides utilities for Gmail API operations.

## Features

- **Send emails** - Simple and advanced message composition
- **Read messages** - List, filter, and retrieve emails
- **Batch operations** - Delete multiple messages efficiently
- **Mail merge** - Send templated emails using Google Sheets data
- **Label management** - List and manage Gmail labels

## Quick Start

```go
import (
    "context"
    "net/http"
    "github.com/grokify/gogoogle/gmailutil/v1"
)

// Create Gmail service
func getGmailService(ctx context.Context, httpClient *http.Client) (*gmailutil.GmailService, error) {
    return gmailutil.NewGmailService(ctx, httpClient)
}
```

## Sending Emails

### SendSimple (Recommended)

The simplest way to send emails:

```go
_, err := service.SendSimple(ctx, "me", gmailutil.SendSimpleOpts{
    To:       "recipient@example.com",
    Subject:  "Hello",
    BodyText: "Plain text body",
    BodyHTML: "<p>HTML body</p>",  // Optional
    ReplyTo:  "replies@example.com", // Optional
})
```

### Send with MessageWriter

For more control over email composition:

```go
import "github.com/grokify/mogo/net/mailutil"

msg := mailutil.MessageWriter{
    To:      mailutil.Addresses{{Address: "recipient@example.com"}},
    Subject: "Hello",
    // ... additional fields
}

_, err := service.Send(ctx, "me", msg)
```

## Reading Messages

### List Messages

```go
messages, err := service.MessagesAPI.GetMessagesList(gmailutil.MessagesListOpts{
    UserID:   "me",
    LabelIDs: []string{"INBOX"},
    MaxResults: 100,
})
```

### Filter by Sender

```go
messages, err := service.MessagesAPI.GetMessagesFrom("sender@example.com")
```

### Filter by Category

```go
// Categories: FORUMS, PROMOTIONS, SOCIAL, UPDATES
messages, err := service.MessagesAPI.GetMessagesByCategory("me", "PROMOTIONS", true)
```

## Batch Operations

### Delete Messages

```go
messageIDs := []string{"msg1", "msg2", "msg3"}
err := service.MessagesAPI.BatchDeleteMessages("me", messageIDs)
```

### Delete by Sender

```go
senders := []string{"spam@example.com", "unwanted@example.com"}
err := service.MessagesAPI.DeleteMessagesFrom(senders)
```

## OAuth Scopes

| Scope | Constant | Description |
|-------|----------|-------------|
| Send only | `GmailSendScope` | Send emails |
| Read only | `GmailReadonlyScope` | Read emails |
| Modify | `GmailModifyScope` | Read, send, delete |
| Full access | `GmailScope` | All operations |

```go
import "github.com/grokify/gogoogle/gmailutil/v1"

// For sending only
scopes := []string{gmailutil.GmailSendScope}

// For full access
scopes := []string{gmailutil.GmailScope}
```

## Next Steps

- [Sending Emails](sending.md) - Detailed sending guide
- [Reading Messages](messages.md) - Query and filter messages
- [Mail Merge](mail-merge.md) - Template-based campaigns
