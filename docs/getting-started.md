# Getting Started

## Installation

```bash
go get github.com/grokify/gogoogle
```

## Authentication

GoGoogle uses OAuth2 for authentication. You'll need:

1. A Google Cloud project with APIs enabled
2. OAuth2 credentials (client_secret.json)
3. User authorization (generates token.json)

### Enable Google APIs

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create or select a project
3. Enable the APIs you need:
   - Gmail API
   - Google Sheets API
   - Google Slides API
   - etc.

### Create OAuth Credentials

1. Go to **APIs & Services > Credentials**
2. Click **Create Credentials > OAuth client ID**
3. Select **Desktop app** as application type
4. Download the JSON file as `client_secret.json`

### Authorize Your Application

Use [goauth](https://github.com/grokify/goauth) for OAuth2 flow:

```go
import (
    "context"
    "github.com/grokify/goauth/google"
)

client, err := google.NewClientOAuthCLITokenStore(google.ClientOAuthCLITokenStoreConfig{
    Context:   ctx,
    AppConfig: clientSecretJSON,  // Contents of client_secret.json
    Scopes:    []string{
        "https://www.googleapis.com/auth/gmail.send",
        "https://www.googleapis.com/auth/spreadsheets.readonly",
    },
    TokenFile: "token.json",      // Will be created after authorization
    State:     "my-app",
})
if err != nil {
    log.Fatal(err)
}
```

On first run, this opens a browser for user authorization. The token is saved to `token.json` for subsequent runs.

## Common OAuth Scopes

### Gmail

| Scope | Description |
|-------|-------------|
| `gmail.send` | Send emails only |
| `gmail.readonly` | Read emails only |
| `gmail.modify` | Read, send, delete emails |
| `gmail.labels` | Manage labels |

```go
import "github.com/grokify/gogoogle/gmailutil/v1"

scopes := []string{gmailutil.GmailSendScope}
```

### Sheets

| Scope | Description |
|-------|-------------|
| `spreadsheets.readonly` | Read spreadsheets |
| `spreadsheets` | Read and write spreadsheets |

### Slides

| Scope | Description |
|-------|-------------|
| `presentations.readonly` | Read presentations |
| `presentations` | Read and write presentations |

## Create Service Clients

After authentication, create service clients:

### Gmail

```go
import "github.com/grokify/gogoogle/gmailutil/v1"

service, err := gmailutil.NewGmailService(ctx, httpClient)
if err != nil {
    log.Fatal(err)
}
```

### Sheets

```go
import "google.golang.org/api/sheets/v4"

service, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
if err != nil {
    log.Fatal(err)
}
```

### Slides

```go
import "google.golang.org/api/slides/v1"

service, err := slides.NewService(ctx, option.WithHTTPClient(httpClient))
if err != nil {
    log.Fatal(err)
}
```

## Environment Variables

For production, store credentials securely:

```bash
export GOOGLE_CLIENT_SECRET=$(cat client_secret.json)
export GOOGLE_TOKEN=$(cat token.json)
```

```go
clientSecret := []byte(os.Getenv("GOOGLE_CLIENT_SECRET"))
```

## Next Steps

- [Gmail Guide](gmail/index.md) - Send emails and manage messages
- [Sheets Guide](sheets/index.md) - Read and write spreadsheet data
- [CLI Tools](cli/index.md) - Command-line utilities
