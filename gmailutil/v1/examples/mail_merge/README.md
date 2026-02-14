# Mail Merge Example

This example demonstrates sending personalized newsletter emails using Mustache templates and recipient data from a Google Sheet.

## Features

- Mustache templates for subject and body (HTML + plain text)
- Inline images via Content-ID references
- File attachments
- Recipient data from Google Sheets

## Prerequisites

1. A Google Cloud project with Gmail and Sheets APIs enabled
2. OAuth2 credentials configured via [goauth](https://github.com/grokify/goauth)
3. A Google Sheet with recipient data

## Google Sheet Setup

Create a Google Sheet with the following columns:

| Column | Description |
|--------|-------------|
| `TO` | Recipient email address (required) |
| `CC` | CC email addresses (optional) |
| `BCC` | BCC email addresses (optional) |
| `FIRST_NAME` | Recipient's first name |
| `COMPANY_NAME` | Your company name |
| `MONTH` | Newsletter month (e.g., "February 2026") |
| `HEADLINE` | Featured article headline |
| `SUMMARY` | Brief summary of the article |
| `ARTICLE_URL` | Link to full article |
| `UNSUBSCRIBE_URL` | Unsubscribe link |

## Command Line Options

| Flag | Short | Description |
|------|-------|-------------|
| `--goauth-credentials-file` | `-c` | Path to goauth credentials JSON file (required) |
| `--goauth-account-key` | `-k` | Account key within credentials file |
| `--sheet-id` | `-s` | Google Sheet ID with recipient data |
| `--sheet-index` | `-x` | Sheet index within spreadsheet (default: 0) |
| `--sheet-header-row-count` | `-r` | Number of header rows (default: 1) |
| `--subject-template` | `-j` | Subject template file |
| `--html-template` | | HTML body template file |
| `--text-template` | `-t` | Plain text body template file |
| `--inline-filename` | `-i` | Inline file (can be repeated) |
| `--attachment-filename` | `-a` | Attachment file (can be repeated) |

## Usage

```bash
go run main.go \
  --goauth-credentials-file=/path/to/credentials.json \
  --goauth-account-key=myaccount \
  --sheet-id=YOUR_GOOGLE_SHEET_ID \
  --subject-template=subject.mustache \
  --html-template=body_html.mustache \
  --text-template=body_text.mustache \
  --inline-filename=logo.png
```

## Template Files

- `subject.mustache` - Email subject line template
- `body_html.mustache` - HTML email body with CSS styling and inline logo
- `body_text.mustache` - Plain text fallback for email clients without HTML support
- `logo.png` - Placeholder logo image (replace with your own)

## Inline Images

To include images in the HTML body, use Content-ID references:

```html
<img src="cid:logo.png" alt="Company Logo" />
```

Then include the image file with `--inline-filename=logo.png`.

## Using with gogoogle CLI

Alternatively, use the unified `gogoogle` CLI:

```bash
gogoogle gmail merge \
  --goauth-credentials-file=/path/to/credentials.json \
  --goauth-credentials-account=myaccount \
  --sheet-id=YOUR_GOOGLE_SHEET_ID \
  --subject-template=subject.mustache \
  --html-template=body_html.mustache \
  --text-template=body_text.mustache \
  --inline=logo.png
```
