# CLI Tools

The `gogoogle` command provides CLI tools for common Google API operations.

## Installation

```bash
go install github.com/grokify/gogoogle/cmd/gogoogle@latest
```

## Commands

| Command | Description |
|---------|-------------|
| `gmail merge` | Send templated emails via mail merge |
| `gmail send-markdown` | Send email with markdown body |
| `slides content` | Extract content from presentations |

## Gmail: Mail Merge

Send templated emails using Google Sheets data:

```bash
gogoogle gmail merge \
    --credentials client_secret.json \
    --token token.json \
    --sheet-id "1abc123..." \
    --subject-template templates/subject.mustache \
    --html-template templates/body.html.mustache \
    --text-template templates/body.txt.mustache
```

### Options

| Flag | Description |
|------|-------------|
| `--credentials` | OAuth credentials file |
| `--token` | Token file path |
| `--sheet-id` | Google Sheets ID |
| `--sheet-name` | Sheet name (default: first sheet) |
| `--subject-template` | Subject template file |
| `--html-template` | HTML body template file |
| `--text-template` | Plain text template file |
| `--inline` | Inline image (format: `cid:path`) |
| `--attachment` | File attachment path |
| `--from` | From address (default: "me") |

### Example with Inline Images

```bash
gogoogle gmail merge \
    --sheet-id "1abc123..." \
    --subject-template "Newsletter - {{date}}" \
    --html-template templates/newsletter.html.mustache \
    --inline "logo:images/logo.png" \
    --inline "banner:images/banner.jpg" \
    --attachment "documents/terms.pdf"
```

## Gmail: Send Markdown

Send a single email with markdown body:

```bash
gogoogle gmail send-markdown \
    --credentials client_secret.json \
    --token token.json \
    --to recipient@example.com \
    --subject "Meeting Notes" \
    --body @notes.md
```

### Options

| Flag | Description |
|------|-------------|
| `--credentials` | OAuth credentials file |
| `--token` | Token file path |
| `--from` | From address (default: "me") |
| `--to` | Recipient address |
| `--cc` | CC recipients (comma-separated) |
| `--bcc` | BCC recipients (comma-separated) |
| `--subject` | Email subject |
| `--body` | Body text or @filename |

### Body from File

Use `@` prefix to read from file:

```bash
gogoogle gmail send-markdown \
    --to user@example.com \
    --subject "Report" \
    --body @report.md
```

### Body Inline

```bash
gogoogle gmail send-markdown \
    --to user@example.com \
    --subject "Quick Note" \
    --body "# Hello\n\nThis is a **quick** note."
```

## Slides: Extract Content

Extract text, images, and notes from a presentation:

```bash
gogoogle slides content \
    --credentials client_secret.json \
    --token token.json \
    --presentation-id "1xyz789..." \
    --output content.json
```

### Options

| Flag | Description |
|------|-------------|
| `--credentials` | OAuth credentials file |
| `--token` | Token file path |
| `--presentation-id` | Google Slides presentation ID |
| `--output` | Output JSON file |
| `--format` | Output format (json, text) |

### Output Format

```json
{
  "presentationId": "1xyz789...",
  "title": "My Presentation",
  "slides": [
    {
      "slideId": "slide1",
      "title": "Introduction",
      "textContent": "Welcome to the presentation...",
      "notes": "Remember to greet the audience",
      "images": [
        {
          "url": "https://...",
          "contentUrl": "https://..."
        }
      ]
    }
  ]
}
```

### Text-Only Output

```bash
gogoogle slides content \
    --presentation-id "1xyz789..." \
    --format text
```

## Authentication

All commands require OAuth authentication:

1. **First run** - Opens browser for authorization
2. **Subsequent runs** - Uses saved token

### Credentials File

Download from Google Cloud Console:

1. Go to **APIs & Services > Credentials**
2. Create **OAuth client ID** (Desktop app)
3. Download as `client_secret.json`

### Token Storage

Tokens are saved to the path specified by `--token`:

```bash
# Default location
~/.gogoogle/token.json

# Custom location
gogoogle gmail merge --token /secure/path/token.json ...
```

## Environment Variables

```bash
export GOOGLE_APPLICATION_CREDENTIALS=client_secret.json
export GOGOOGLE_TOKEN_FILE=~/.gogoogle/token.json
```

## Examples

### Weekly Newsletter

```bash
#!/bin/bash
# Send weekly newsletter

gogoogle gmail merge \
    --sheet-id "$NEWSLETTER_SHEET_ID" \
    --subject-template "Weekly Update - $(date +%B\ %d)" \
    --html-template templates/weekly.html.mustache \
    --text-template templates/weekly.txt.mustache \
    --inline "logo:assets/logo.png"
```

### Export Presentation to Text

```bash
#!/bin/bash
# Export all presentations in a folder

for id in $(cat presentation_ids.txt); do
    gogoogle slides content \
        --presentation-id "$id" \
        --output "exports/${id}.json"
done
```

## Next Steps

- [Gmail Guide](../gmail/index.md) - Detailed Gmail documentation
- [Slides Guide](../slides/index.md) - Slides API documentation
