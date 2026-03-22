# Reading Messages

## List Messages

### Basic Listing

```go
import "github.com/grokify/gogoogle/gmailutil/v1"

messages, err := service.MessagesAPI.GetMessagesList(gmailutil.MessagesListOpts{
    UserID:     "me",
    MaxResults: 100,
})
```

### With Labels

```go
messages, err := service.MessagesAPI.GetMessagesList(gmailutil.MessagesListOpts{
    UserID:   "me",
    LabelIDs: []string{"INBOX", "UNREAD"},
})
```

### With Query

Use Gmail search syntax:

```go
messages, err := service.MessagesAPI.GetMessagesList(gmailutil.MessagesListOpts{
    UserID: "me",
    Query:  "from:important@example.com after:2024/01/01",
})
```

## MessagesListOpts

| Field | Type | Description |
|-------|------|-------------|
| `UserID` | `string` | User ID ("me" for authenticated user) |
| `LabelIDs` | `[]string` | Filter by labels |
| `Query` | `string` | Gmail search query |
| `MaxResults` | `int64` | Maximum messages to return |
| `PageToken` | `string` | Pagination token |
| `IncludeSpamTrash` | `bool` | Include spam/trash |

## Query Syntax

Gmail search operators:

| Operator | Example | Description |
|----------|---------|-------------|
| `from:` | `from:user@example.com` | Messages from sender |
| `to:` | `to:me@example.com` | Messages to recipient |
| `subject:` | `subject:meeting` | Subject contains |
| `after:` | `after:2024/01/01` | After date |
| `before:` | `before:2024/12/31` | Before date |
| `is:` | `is:unread` | Message state |
| `has:` | `has:attachment` | Has attachment |
| `category:` | `category:promotions` | Gmail category |

### Query Examples

```go
// Unread from specific sender
query := "from:boss@company.com is:unread"

// Messages with attachments this year
query := "has:attachment after:2024/01/01"

// In a specific category
query := "category:updates"

// Combine multiple criteria
query := "from:notifications@github.com subject:PR after:2024/06/01"
```

## Filter by Sender

```go
// Get messages from specific sender
messages, err := service.MessagesAPI.GetMessagesFrom("sender@example.com")
```

## Filter by Category

Gmail categories: `FORUMS`, `PROMOTIONS`, `SOCIAL`, `UPDATES`

```go
// Get all promotional emails
messages, err := service.MessagesAPI.GetMessagesByCategory("me", "PROMOTIONS", true)

// Get first page only
messages, err := service.MessagesAPI.GetMessagesByCategory("me", "SOCIAL", false)
```

## Get Full Message

Message list returns metadata only. Get full content:

```go
// Get message with full content
message, err := service.MessagesAPI.GetMessage("me", messageID)

// Access headers
for _, header := range message.Payload.Headers {
    if header.Name == "Subject" {
        fmt.Println("Subject:", header.Value)
    }
}
```

## Inflate Messages

Convert message metadata to full messages:

```go
// Get list (metadata only)
list, _ := service.MessagesAPI.GetMessagesList(opts)

// Inflate to full messages
fullMessages, err := service.MessagesAPI.InflateMessages("me", list.Messages)
```

## Batch Delete

### Delete by IDs

```go
messageIDs := []string{"msg1", "msg2", "msg3"}
err := service.MessagesAPI.BatchDeleteMessages("me", messageIDs)
```

### Delete by Sender

Delete all messages from specific senders:

```go
senders := []string{
    "spam@example.com",
    "unwanted@example.com",
}
err := service.MessagesAPI.DeleteMessagesFrom(senders)
```

## Pagination

Handle large result sets:

```go
var allMessages []*gmail.Message
pageToken := ""

for {
    result, err := service.MessagesAPI.GetMessagesList(gmailutil.MessagesListOpts{
        UserID:     "me",
        MaxResults: 100,
        PageToken:  pageToken,
    })
    if err != nil {
        return err
    }

    allMessages = append(allMessages, result.Messages...)

    if result.NextPageToken == "" {
        break
    }
    pageToken = result.NextPageToken
}
```

## Labels

### List Labels

```go
labels, err := gmailutil.GetLabelNames(httpClient)
for _, label := range labels {
    fmt.Println(label)
}
```

### Common Label IDs

| Label | ID |
|-------|-----|
| Inbox | `INBOX` |
| Sent | `SENT` |
| Drafts | `DRAFT` |
| Spam | `SPAM` |
| Trash | `TRASH` |
| Starred | `STARRED` |
| Unread | `UNREAD` |
| Important | `IMPORTANT` |

## Next Steps

- [Sending Emails](sending.md) - Send messages
- [Mail Merge](mail-merge.md) - Bulk sending with templates
