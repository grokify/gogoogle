# Google Drive

The `driveutil/v3` package uploads local files to Drive, optionally converting
them to Google Docs, Sheets, or Slides, replaces the content of an existing
file, and lists a file's comments. Together with the suggestion-aware Docs
extraction and the Sheets records reader, it supports a review round trip:
publish a document, let people comment and suggest edits, and read their
responses back as data.

## Installation

```bash
go get github.com/grokify/gogoogle/driveutil/v3
```

## Create the service

```go
import driveutil "github.com/grokify/gogoogle/driveutil/v3"

svc, err := driveutil.NewService(ctx, httpClient)
if err != nil {
    return err
}
```

`httpClient` is any authorized `*http.Client`. Extra `option.ClientOption`
values, such as `option.WithEndpoint`, follow the client.

## Upload with conversion

```go
f, err := svc.Upload(ctx, driveutil.UploadRequest{
    Path:      "review.docx",
    FolderID:  folderID, // optional
    ConvertTo: driveutil.ConvertDoc,
})
// f.ID, f.Name, f.MimeType, f.WebViewLink, f.ModifiedTime, f.HeadRevisionID
```

| `ConvertTo` | Accepted extensions | Result |
|-------------|---------------------|--------|
| `ConvertDoc` | `.docx`, `.odt`, `.txt`, `.html` | Google Doc |
| `ConvertSheet` | `.csv`, `.xlsx`, `.ods` | Google Sheet |
| `ConvertSlides` | `.pptx`, `.odp` | Google Slides file |
| `ConvertNone` (empty) | any | stored unchanged |

An extension that does not match the target returns an error wrapping
`ErrUnsupportedConversion`. When converting, the file name defaults to the base
name without its extension.

`sheetsutil.CreateFromCSV(ctx, svc, "data.csv", "", folderID)` is a shortcut
for a CSV-to-Sheet upload.

## Update in place

```go
f, err := svc.Update(ctx, fileID, "review-v2.docx", driveutil.ConvertDoc)
```

The file ID, link, sharing, and comments are kept, so reviewers keep working in
the same document.

## Comments

```go
comments, err := svc.ListComments(ctx, fileID, driveutil.ListCommentsOptions{})
for _, c := range comments {
    fmt.Println(c.QuotedText, c.Content, c.Resolved, len(c.Replies))
}
```

Each comment has the text it is anchored to (`QuotedText`), its author, its
resolved state, and its replies. Paging is followed to the end: you get every
comment or an error, never a partial list. Author email addresses appear only
when the API returns them; treat them as personal data.

## Suggested edits

Comments do not include suggested edits. Fetch the document with its
suggestions shown inline and split the text:

```go
doc, err := docsSvc.GetDocumentWithSuggestions(ctx, docID, "SUGGESTIONS_INLINE")
sc := docsutil.ExtractSuggestionContent(doc)
for _, b := range sc.Blocks {
    fmt.Println(b.Location, b.Original, "->", b.Accepted)
}
```

See the [Docs guide](../docs/index.md#suggestions).

## Structured responses in a Sheet

For answers that need a fixed shape, publish a CSV as a Sheet and read it back:

```go
records, err := sheetsutil.ReadRecords(ctx, sheetsSvc, sheetID, "") // first sheet
// []map[string]string keyed by the header row
```

See the [Sheets guide](../sheets/index.md#records).

## Command line

```bash
gogoogle drive upload review.docx --convert doc --folder 1AbC
gogoogle drive update 1FileId review-v2.docx --convert doc
gogoogle drive comments 1FileId --include-deleted
gogoogle docs suggestions 1DocId
gogoogle sheets records 1SheetId --sheet "Round 1"
```

All commands print JSON. Exit status is 0 on success, 1 for an API or input
error, and 2 for a usage error.

## Scopes

Use the narrowest scope that works.

| Operation | Scope |
|-----------|-------|
| Upload, update | `drive.file` (files this app creates or opens) |
| List comments | `drive.readonly` |
| Read suggestions | `documents.readonly` |
| Read records | `spreadsheets.readonly` |

## Errors

Failures wrap the operation and the resource ID. `ErrNotFound` (HTTP 404),
`ErrPermissionDenied` (HTTP 403), and `ErrUnsupportedConversion` work with
`errors.Is`; the underlying `*googleapi.Error` works with `errors.As`. Retries
are left to the caller's HTTP client.
