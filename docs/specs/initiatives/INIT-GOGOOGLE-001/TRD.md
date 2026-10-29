# TRD — Document review round trip — Drive upload and conversion, comments, and suggestion-aware Docs content

**Initiative:** `INIT-GOGOOGLE-001`
**Repository:** `github.com/grokify/gogoogle`
**Version:** 0.1
**Status:** Draft

MUST / MUST NOT / SHOULD / MAY are used in the RFC 2119 sense.

## 1. Introduction

**Purpose.** Specify the new `driveutil/v3` package, suggestion-aware
extraction in `docsutil/v1`, record helpers in `sheetsutil/v4`, and CLI
commands for a document review round trip.

**Scope.** Library and CLI only; see PRD non-goals.

**References.** [PRD](PRD.md) · [PLAN](PLAN.md) · [ROADMAP](ROADMAP.md) ·
Google Drive API v3 (`files.create`, `files.update`, `files.get`,
`comments.list`) · Google Docs API v1 (`documents.get`,
`suggestionsViewMode`).

## 2. Technical Requirements and Traceability

| TRD ID | Requirement | Source | Verification |
|--------|-------------|--------|--------------|
| TRD-001 | `driveutil/v3.NewService` and `GetFile` | PRD-FR-1 | Fake-server test |
| TRD-002 | `Upload` with `ConvertTo` mapping to Google MIME types; optional `FolderID` | PRD-FR-2 | Request-capture test per conversion |
| TRD-003 | `Update` media-replaces an existing file ID | PRD-FR-3 | Request-capture test |
| TRD-004 | `ListComments` pages through all results with the fields in §6.3 | PRD-FR-4 | Multi-page fake-server test |
| TRD-005 | `docsutil.ExtractSuggestionContent` per §6.4 | PRD-FR-5 | Recorded-response fixtures |
| TRD-006 | `sheetsutil.ReadRecords` and `CreateFromCSV` | PRD-FR-6 | Fake-server test |
| TRD-007 | CLI commands in §6.5 | PRD-FR-7 | CLI tests |
| TRD-008 | No exported identifier removed or changed in signature | PRD-FR-8 | `apidiff` against v0.12.0 in CI |

## 3. System Invariants

- **INV-1** New functions never mutate caller inputs; all return new values.
- **INV-2** Every API call takes a `context.Context`; errors are wrapped with
  the operation and resource ID and returned, never logged-and-dropped.
- **INV-3** Listing functions return either all pages or an error — never a
  silent partial result.
- **INV-4** Suggestion extraction preserves document order; concatenating
  `original` + `deletion` segments yields the original text, and `original` +
  `insertion` segments yields the as-if-accepted text.

## 4. Architecture

```
caller ─► driveutil/v3  ──► Drive API v3 (files, comments)
       ─► docsutil/v1   ──► Docs API v1 (documents.get, SUGGESTIONS_INLINE)
       ─► sheetsutil/v4 ──► Sheets API v4 (values.get); CSV → Sheet via driveutil
cmd/gogoogle (Cobra) ─► the packages above, JSON output
```

| Component | Responsibilities |
|-----------|------------------|
| `driveutil/v3` (new) | Service, file get, upload / update with conversion, comment listing |
| `docsutil/v1` (extended) | Suggestion-aware extraction alongside existing `ExtractDocumentContent` |
| `sheetsutil/v4` (extended) | `ReadRecords`; `CreateFromCSV` delegating to `driveutil` |
| `cmd/gogoogle` | `drive`, `docs suggestions`, `sheets records` subcommands |

## 5. State and Lifecycle

Not applicable — stateless client wrappers; Drive owns file and comment state.

## 6. Interface Contracts

### 6.1 Drive service and files

```go
package driveutil // import "github.com/grokify/gogoogle/driveutil/v3"

type Service struct{ /* wraps *drive.Service */ }
func NewService(ctx context.Context, c *http.Client, opts ...option.ClientOption) (*Service, error)

type File struct {
    ID, Name, MimeType, WebViewLink, HeadRevisionID string
    ModifiedTime time.Time
}
func (s *Service) GetFile(ctx context.Context, fileID string) (*File, error)
```

### 6.2 Upload and update

```go
type ConvertTo string // "", "doc", "sheet", "slides"
type UploadRequest struct {
    Path      string    // local file; media type inferred from extension
    Name      string    // defaults to base name without extension when converting
    FolderID  string    // optional parent
    ConvertTo ConvertTo
}
func (s *Service) Upload(ctx context.Context, r UploadRequest) (*File, error)
func (s *Service) Update(ctx context.Context, fileID, path string, convertTo ConvertTo) (*File, error)
```

**Conversion targets:**

| `ConvertTo` | Google MIME type | Allowed source extensions |
|-------------|------------------|---------------------------|
| `doc` | `application/vnd.google-apps.document` | `.docx`, `.odt`, `.txt`, `.html` |
| `sheet` | `application/vnd.google-apps.spreadsheet` | `.csv`, `.xlsx`, `.ods` |
| `slides` | `application/vnd.google-apps.presentation` | `.pptx`, `.odp` |

Any other pair returns `ErrUnsupportedConversion`.

### 6.3 Comments

```go
type Comment struct {
    ID, Content, QuotedText, Anchor string
    Author                          Author
    CreatedTime, ModifiedTime       time.Time
    Resolved, Deleted               bool
    Replies                         []Reply
}
type Reply struct {
    ID, Content, Action string // action: "", "resolve", "reopen"
    Author              Author
    CreatedTime         time.Time
}
type Author struct{ DisplayName, EmailAddress string } // email only when the API returns it
type ListCommentsOptions struct {
    IncludeDeleted    bool
    StartModifiedTime time.Time
}
func (s *Service) ListComments(ctx context.Context, fileID string, o ListCommentsOptions) ([]Comment, error)
```

Uses `comments.list` with an explicit `fields` mask and `pageSize=100`,
following `nextPageToken` to the end (INV-3).

### 6.4 Suggestion-aware Docs extraction

```go
type SegmentKind string // "original", "insertion", "deletion"
type Segment struct {
    Kind          SegmentKind `json:"kind"`
    Text          string      `json:"text"`
    SuggestionIDs []string    `json:"suggestion_ids,omitempty"`
}
type Block struct {
    Location string    `json:"location"` // "body/p/12" or "body/table/2/r/3/c/1"
    Segments []Segment `json:"segments"`
    Original string    `json:"original"` // original + deletion segments
    Accepted string    `json:"accepted"` // original + insertion segments
}
type SuggestionContent struct {
    DocumentID    string   `json:"document_id"`
    Blocks        []Block  `json:"blocks"`
    SuggestionIDs []string `json:"suggestion_ids"` // distinct, document order
}
func ExtractSuggestionContent(doc *docs.Document) *SuggestionContent
```

**Fetching.** Callers fetch the document with
`GetDocumentWithSuggestions(ctx, id, "SUGGESTIONS_INLINE")`.

**Classification:**

- A text run with `SuggestedInsertionIds` → `insertion`.
- A text run with `SuggestedDeletionIds` → `deletion`.
- Any other run → `original`.

**Grouping.** Adjacent runs of the same kind and IDs are merged. Table cells
are walked in row-major order. A block is emitted for every paragraph and
cell, including blocks without suggestions, so callers can align blocks with
`ExtractDocumentContent` output.

**Authors.** The Docs API does not expose suggestion authors, so `Block`
carries none. Comments do carry authors.

### 6.5 CLI

| Command | Output |
|---------|--------|
| `gogoogle drive upload <path> [--name] [--folder ID] [--convert doc\|sheet\|slides]` | `File` JSON |
| `gogoogle drive update <fileID> <path> [--convert …]` | `File` JSON |
| `gogoogle drive comments <fileID> [--include-deleted]` | `[]Comment` JSON |
| `gogoogle docs suggestions <documentID>` | `SuggestionContent` JSON |
| `gogoogle sheets records <spreadsheetID> [--sheet NAME]` | `[]map[string]string` JSON |

Exit codes: 0 success; 1 API or input error; 2 usage error.

### 6.6 Errors

**Sentinels.** These are wrapped with context:

- `ErrUnsupportedConversion`
- `ErrNotFound` (HTTP 404)
- `ErrPermissionDenied` (HTTP 403)

**Other errors** wrap the `googleapi.Error`.

**Retries** are left to the caller's HTTP client.

## 7. Data Design

No persistence. JSON tags use snake_case, matching `docsutil`'s existing
convention, across the new types.

## 8. Operational Guarantees

| Category | Requirement | Measurement | Applicability |
|----------|-------------|-------------|---------------|
| Performance | One API round trip per page; replies come in the same `comments.list` response (no N+1) | Request-count assertion in tests | Applicable |
| Reliability | INV-3 complete-or-error listings | Paged tests | Applicable |
| Availability, throughput, capacity | Not applicable — client library | — | — |
| Observability | Errors carry operation and resource ID | Unit tests | Applicable |

## 9. Security and Privacy

- **Scopes:**
  - Reading uses `drive.readonly`.
  - Upload and update need `drive.file`, which is limited to files the app
    creates or opens.
  - The guide documents the narrowest scope for each function.
- **Data handling:**
  - No credentials or tokens are logged.
  - Comment author emails are returned only when the API provides them. The
    guide advises treating them as personal data.
- **Encryption:** Not applicable to this library. The Google client uses TLS
  in transit, and the library adds no storage.

## 10. Compatibility and Migration

Applicable: additive only within v0.x minor (TRD-008). No migration.

## 11. Dependencies

| Dependency | Version | Purpose |
|------------|---------|---------|
| `google.golang.org/api` (`drive/v3`, `docs/v1`, `sheets/v4`) | current in go.mod (v0.307.0); bump to latest verified version at implementation | Google APIs |
| `github.com/spf13/cobra` | v1.10.2 | CLI |

No new modules: `drive/v3` is part of `google.golang.org/api`.

## 12. Testing Strategy

- **Fake servers.** Unit tests use `httptest` servers with
  `option.WithEndpoint` and `option.WithHTTPClient`, recording requests to
  assert methods, MIME types, `fields` masks, and paging.
- **Suggestion fixtures.** Suggestion extraction is tested against recorded
  `documents.get` JSON responses, plus a property test for INV-4.
- **Optional integration tests.** These run behind `GOGOOGLE_IT=1` and are
  never required in CI.

## 13. Deployment

Released as a tagged minor version (v0.13.0) after CI is green; release
notes list the new packages and functions. Rollback: callers pin v0.12.0.

## 14. Alternatives Considered

| Alternative | Why rejected |
|-------------|--------------|
| Read comments by exporting the file as `.docx` and parsing comment XML | Loses resolved state and reply threading and needs a second parser; the Drive comments API is authoritative |
| Return suggestions as raw `docs.TextRun` lists | Leaves every caller to reimplement classification; the segment model with original and accepted views is what consumers need |
| Implement these features inside MCP wrappers | Wrappers should stay thin; other Go programs and the CLI need the same functions |

## 15. Risks and Open Questions

| ID | Risk / Question | Mitigation |
|----|-----------------|------------|
| R-1 | Drive conversion fidelity for `.docx` / `.pptx` (tables, images) varies | Documented; callers keep source formats simple |
| R-2 | Suggestion authors unavailable in the Docs API | Documented; comments carry authors |
| R-3 | Comment anchors for converted files are opaque strings | Return `QuotedText` as the primary anchor; `Anchor` passed through unparsed |

No open question blocks implementation.

## Appendix: Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1 | 2026-10-29 | John Wang | Initial draft |
