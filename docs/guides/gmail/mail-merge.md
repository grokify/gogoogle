# Mail Merge

Send templated emails to multiple recipients using Google Sheets as the data source.

## Overview

Mail merge allows you to:

- Define email templates with Mustache placeholders
- Pull recipient data from Google Sheets
- Send personalized emails at scale
- Include inline images and attachments

## Quick Start

```go
import "github.com/grokify/gogoogle/gmailutil/v1/mailmerge"

mm := mailmerge.MailMerge{
    GmailService: gmailService,
    SheetsService: sheetsService,
    SpreadsheetID: "your-spreadsheet-id",
    SubjectTemplate: "Hello {{name}}!",
    BodyHTMLTemplate: "<p>Dear {{name}}, your order #{{order_id}} is ready.</p>",
    BodyTextTemplate: "Dear {{name}}, your order #{{order_id}} is ready.",
}

err := mm.Execute(ctx)
```

## Google Sheets Format

Your spreadsheet must have these columns:

| Column | Required | Description |
|--------|----------|-------------|
| `TO` | Yes | Recipient email address |
| `CC` | No | CC recipients |
| `BCC` | No | BCC recipients |
| `*` | No | Any other columns become template variables |

### Example Sheet

| TO | CC | name | order_id | amount |
|----|-----|------|----------|--------|
| john@example.com | | John | 12345 | $99.00 |
| jane@example.com | manager@example.com | Jane | 12346 | $149.00 |

## Mustache Templates

Templates use [Mustache](https://mustache.github.io/) syntax:

### Subject Template

```
Order Confirmation - {{order_id}}
```

### HTML Body Template

```html
<html>
<body>
    <h1>Hello {{name}}!</h1>
    <p>Your order <b>#{{order_id}}</b> for {{amount}} has been confirmed.</p>
    <p>Thank you for your purchase!</p>
</body>
</html>
```

### Plain Text Body Template

```
Hello {{name}}!

Your order #{{order_id}} for {{amount}} has been confirmed.

Thank you for your purchase!
```

## Template Files

Load templates from files:

```go
mm := mailmerge.MailMerge{
    GmailService:       gmailService,
    SheetsService:      sheetsService,
    SpreadsheetID:      spreadsheetID,
    SubjectTemplateFile: "templates/subject.mustache",
    BodyHTMLTemplateFile: "templates/body.html.mustache",
    BodyTextTemplateFile: "templates/body.txt.mustache",
}
```

## Inline Images

Include images with Content-ID references:

### Template

```html
<img src="cid:logo" alt="Company Logo">
<p>Hello {{name}}!</p>
```

### Configuration

```go
mm := mailmerge.MailMerge{
    // ... other config ...
    InlineFiles: []mailmerge.InlineFile{
        {
            ContentID: "logo",
            FilePath:  "images/logo.png",
        },
    },
}
```

## Attachments

Add file attachments:

```go
mm := mailmerge.MailMerge{
    // ... other config ...
    Attachments: []string{
        "documents/terms.pdf",
        "documents/receipt.pdf",
    },
}
```

## CLI Usage

Use the `gogoogle` CLI for mail merge:

```bash
gogoogle gmail merge \
    --sheet-id "1234567890abcdef" \
    --subject-template "templates/subject.mustache" \
    --html-template "templates/body.html.mustache" \
    --text-template "templates/body.txt.mustache" \
    --inline "logo:images/logo.png" \
    --attachment "documents/terms.pdf"
```

## MailMerge Options

| Field | Type | Description |
|-------|------|-------------|
| `GmailService` | `*gmailutil.GmailService` | Gmail API service |
| `SheetsService` | `*sheets.Service` | Sheets API service |
| `SpreadsheetID` | `string` | Google Sheets ID |
| `SheetName` | `string` | Sheet name (default: first sheet) |
| `SubjectTemplate` | `string` | Subject template string |
| `SubjectTemplateFile` | `string` | Subject template file path |
| `BodyHTMLTemplate` | `string` | HTML body template string |
| `BodyHTMLTemplateFile` | `string` | HTML body template file path |
| `BodyTextTemplate` | `string` | Plain text body template string |
| `BodyTextTemplateFile` | `string` | Plain text body template file path |
| `InlineFiles` | `[]InlineFile` | Inline images with Content-ID |
| `Attachments` | `[]string` | File paths to attach |
| `FromAddress` | `string` | Sender address (default: "me") |

## Example: Newsletter

### Sheet (Newsletter Recipients)

| TO | name | unsubscribe_token |
|----|------|-------------------|
| user1@example.com | Alice | abc123 |
| user2@example.com | Bob | def456 |

### subject.mustache

```
Weekly Newsletter - March 2024
```

### body.html.mustache

```html
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; }
        .header { background: #4285f4; color: white; padding: 20px; }
        .content { padding: 20px; }
        .footer { font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="header">
        <img src="cid:logo" alt="Logo" height="40">
    </div>
    <div class="content">
        <h2>Hello {{name}}!</h2>
        <p>Here's what's new this week...</p>
    </div>
    <div class="footer">
        <a href="https://example.com/unsubscribe?token={{unsubscribe_token}}">Unsubscribe</a>
    </div>
</body>
</html>
```

### Execute

```go
mm := mailmerge.MailMerge{
    GmailService:         gmailService,
    SheetsService:        sheetsService,
    SpreadsheetID:        "newsletter-sheet-id",
    SubjectTemplateFile:  "templates/subject.mustache",
    BodyHTMLTemplateFile: "templates/body.html.mustache",
    BodyTextTemplateFile: "templates/body.txt.mustache",
    InlineFiles: []mailmerge.InlineFile{
        {ContentID: "logo", FilePath: "images/logo.png"},
    },
}

if err := mm.Execute(ctx); err != nil {
    log.Fatal(err)
}
```

## Best Practices

1. **Test first** - Send to yourself before bulk sending
2. **Include unsubscribe** - Required for marketing emails
3. **Plain text fallback** - Always include text version
4. **Rate limiting** - Gmail has daily sending limits
5. **Validate data** - Check sheet for missing/invalid emails

## Next Steps

- [Sending Emails](sending.md) - Basic email sending
- [CLI Tools](../cli/index.md) - Command-line mail merge
