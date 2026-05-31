# GoGoogle

Go libraries for Google APIs - higher-level utilities built on the official Google API Go Client.

## Overview

GoGoogle provides simplified, production-ready interfaces for common Google API operations:

| Package | Description |
|---------|-------------|
| [Gmail](gmail/index.md) | Send emails, manage messages, mail merge |
| [Docs](docs/index.md) | Extract content from Google Docs |
| [Sheets](sheets/index.md) | Read/write spreadsheet data with typed structs |
| [Slides](slides/index.md) | Create presentations, add content |
| [Maps](maps/index.md) | Generate static map images |
| [Speech](speech/stt.md) | Speech-to-Text and Text-to-Speech |
| [Forms](#forms) | OAuth scope helpers for Forms API |
| [YouTube](#youtube) | URL utilities (short URLs) |
| [CLI](cli/index.md) | Command-line tools for Google APIs |

## Features

- **Simplified APIs** - Higher-level functions for common operations
- **Type Safety** - Go structs for API data
- **Batch Operations** - Efficient bulk processing
- **Mail Merge** - Template-based email campaigns with Sheets integration
- **OAuth2 Support** - Built-in authentication helpers

## Quick Start

```bash
go get github.com/grokify/gogoogle
```

### Send an Email

```go
import (
    "context"
    "github.com/grokify/gogoogle/gmailutil/v1"
)

// Create Gmail service (after OAuth setup)
service, _ := gmailutil.NewGmailService(ctx, httpClient)

// Send a simple email
service.SendSimple(ctx, "me", gmailutil.SendSimpleOpts{
    To:       "recipient@example.com",
    Subject:  "Hello from GoGoogle",
    BodyText: "This is a test email.",
    BodyHTML: "<p>This is a <b>test</b> email.</p>",
})
```

### Read Google Sheets

```go
import "github.com/grokify/gogoogle/sheetsutil/v4/sheetsmap"

// Read sheet data into typed structs
data, _ := sheetsmap.ReadSheet(ctx, service, spreadsheetID, "Sheet1")
for _, row := range data.Rows {
    fmt.Println(row["Name"], row["Email"])
}
```

## Forms

The `forms/v1` package provides OAuth scope helpers for the Google Forms API.

```go
import forms "github.com/grokify/gogoogle/forms/v1"

// All Forms scopes (Drive, DriveFile, FormsBody, FormsResponsesReadonly)
scopes := forms.ScopesAll()

// Read-only scopes
scopes := forms.ScopesReadOnly()

// As client option for google.golang.org/api
opt := forms.ClientOptionScopesAll()
```

## YouTube

The `youtubeutil` package provides URL utilities for YouTube.

```go
import "github.com/grokify/gogoogle/youtubeutil"

// Convert YouTube URLs to short format
shortURL, err := youtubeutil.ShortURL("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
// Returns: "https://youtu.be/dQw4w9WgXcQ"

// Also accepts video ID directly
shortURL, err := youtubeutil.ShortURL("dQw4w9WgXcQ")
// Returns: "https://youtu.be/dQw4w9WgXcQ"
```

## Related Libraries

- [goauth](https://github.com/grokify/goauth) - OAuth2 utilities for Google and other providers
- [mogo](https://github.com/grokify/mogo) - General Go utilities used by GoGoogle

## License

MIT License - see [LICENSE](https://github.com/grokify/gogoogle/blob/main/LICENSE) for details.
