# Google Docs Utilities

The `docsutil/v1` package provides utilities for working with the Google Docs API.

## Installation

```go
import docsutil "github.com/grokify/gogoogle/docsutil/v1"
```

## Features

- **Client wrapper** - Simplified Google Docs API client
- **Content extraction** - Extract structured content (headings, paragraphs, images, tables)
- **Text extraction** - Get plain text or paragraph-by-paragraph text
- **URL parsing** - Extract document IDs from Google Docs URLs

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"

    docsutil "github.com/grokify/gogoogle/docsutil/v1"
)

func main() {
    ctx := context.Background()

    // Create client with authenticated HTTP client
    client, err := docsutil.NewClient(ctx, httpClient)
    if err != nil {
        log.Fatal(err)
    }

    // Get document by ID or URL
    docID := docsutil.ParseDocumentID("https://docs.google.com/document/d/DOC_ID/edit")

    // Extract all content
    content, err := client.ExtractContent(ctx, docID, docsutil.ExtractOptions{
        IncludeImages: true,
        IncludeTables: true,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Title: %s\n", content.Title)
    for _, section := range content.Sections {
        fmt.Printf("[%s] %s\n", section.Type, section.Text)
    }
}
```

## Content Extraction

### ExtractContent

Returns structured content including headings, paragraphs, images, and tables.

```go
content, err := client.ExtractContent(ctx, docID, docsutil.ExtractOptions{
    IncludeImages:  true,
    IncludeTables:  true,
    IncludeHeaders: true,
    IncludeFooters: true,
})
```

**Returns:**

- `Title` - Document title
- `Sections` - Array of content sections with type, level, and text
- `Images` - Array of image info (if requested)
- `Tables` - Array of table data (if requested)
- `Headers` - Array of header text (if requested)
- `Footers` - Array of footer text (if requested)

### ExtractText

Returns all document text as a single string.

```go
text, err := client.ExtractText(ctx, docID)
```

### ExtractParagraphs

Returns text organized by paragraphs.

```go
paragraphs, err := client.ExtractParagraphs(ctx, docID)
for _, p := range paragraphs {
    fmt.Println(p)
}
```

## Suggestions

Reviewers' suggested edits do not appear in a normal fetch. Fetch the document
with suggestions shown inline, then split the text into what was there, what
would be added, and what would be removed:

```go
doc, err := svc.GetDocumentWithSuggestions(ctx, docID, "SUGGESTIONS_INLINE")
if err != nil {
    return err
}
sc := docsutil.ExtractSuggestionContent(doc)
for _, b := range sc.Blocks {
    fmt.Printf("%s\n  was: %q\n  if accepted: %q\n", b.Location, b.Original, b.Accepted)
}
```

- Every paragraph and table cell is a block, with a `Location` such as
  `body/p/12` or `body/table/2/r/3/c/1`.
- Each block has segments of kind `original`, `insertion`, or `deletion`, in
  document order, with the suggestion IDs that affect them.
- `Original` is the original and deletion text; `Accepted` is the original and
  insertion text.
- The Docs API does not report who made a suggestion, so blocks carry no
  author. Comments do; see the [Drive guide](../drive/index.md).

From the command line: `gogoogle docs suggestions <documentID|URL>`.

## URL Parsing

Parse document IDs from various Google Docs URL formats:

```go
// Standard edit URL
id := docsutil.ParseDocumentID("https://docs.google.com/document/d/DOC_ID/edit")

// URL with query parameters
id := docsutil.ParseDocumentID("https://docs.google.com/document/d/DOC_ID/edit?tab=t.0")

// URL with anchors
id := docsutil.ParseDocumentID("https://docs.google.com/document/d/DOC_ID/edit#heading=h.xyz")

// Just the ID
id := docsutil.ParseDocumentID("DOC_ID")
```

## Authentication

The client requires an authenticated `*http.Client`. Use [goauth](https://github.com/grokify/goauth) for OAuth2 setup:

```go
import "github.com/grokify/goauth/google"

// Service account
httpClient, err := google.NewClientSvcAccountFromFile(ctx, "credentials.json",
    "https://www.googleapis.com/auth/documents.readonly",
)

// OAuth2 user credentials
httpClient, err := google.NewClientFromTokenFile(ctx, "token.json")
```

## Required Scopes

- `https://www.googleapis.com/auth/documents.readonly` - Read-only access
- `https://www.googleapis.com/auth/documents` - Full access
