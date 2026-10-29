# PRD — Document review round trip — Drive upload and conversion, comments, and suggestion-aware Docs content

**Initiative:** `INIT-GOGOOGLE-001`
**Repository:** `github.com/grokify/gogoogle`
**Workflow:** pbhq-lite
**Status:** Draft — pending stakeholder review (pbhq-lite gate: after PRD)

## 1. Problem

A common workflow is to generate a document or deck from structured data,
share it in Google Drive for people to review, and then read their responses
back into the source system. gogoogle covers parts of this today:

- **Read Docs, Sheets, and Slides.** `docsutil`, `sheetsutil`, and
  `slidesutil` read native Google files.
- **Fetch with suggestions inline.** `docsutil.GetDocumentWithSuggestions`
  can fetch a document in `SUGGESTIONS_INLINE` mode.

Three gaps block the round trip:

1. **No Drive package.** There is no way to upload a `.docx`, `.pptx`, or
   `.csv` *with conversion* to a native Google Doc, Slides, or Sheet, or to
   replace a file's contents with a new version while keeping its ID and
   sharing. Reviewers therefore get files that the Docs, Sheets, and Slides
   APIs cannot read, and every revision becomes a new link.
2. **Comments are invisible.** Reviewers mostly respond with comments.
   Comments live in the Drive API (`comments.list`), not the Docs API, so no
   gogoogle function returns them, along with the text they anchor to, their
   replies, or whether they are resolved.
3. **Suggestions are flattened.** `ExtractDocumentContent` concatenates text
   runs without looking at `suggestedInsertionIds` or
   `suggestedDeletionIds`. A document read in suggestions-inline mode returns
   original text, suggested insertions, and suggested deletions mixed
   together, so a consumer cannot tell what a reviewer proposed.

**Who is affected:**

- Go applications that publish documents for review.
- MCP servers and agents built on gogoogle (they can expose only what
  gogoogle provides).
- CLI users who want a scriptable upload and review read.

**Desired state:** with one library and CLI, a caller can:

- upload with conversion;
- update in place;
- list comments with their anchors;
- read a Doc as segments marked original, suggested insertion, or suggested
  deletion;
- create a Sheet from CSV and read it back as records.

## 2. Goals and Non-Goals

### 2.1 Goals

| ID | Goal | Success criterion |
|----|------|-------------------|
| G-1 | Publish files for review as native Google files | Upload of `.docx` → Doc, `.pptx` → Slides, `.csv`/`.xlsx` → Sheet returns file ID and web link; tests cover each conversion's request |
| G-2 | Revise without breaking links | Update-in-place replaces content of an existing file ID; ID and permissions unchanged |
| G-3 | Read comments with context | Every comment returns author display name, created/modified time, content, quoted anchor text, resolved flag, and replies; all pages fetched |
| G-4 | Read suggestions faithfully | Suggestion-aware extraction distinguishes original text, suggested insertions, and suggested deletions (with suggestion IDs), for paragraphs and table cells |
| G-5 | Review sheets as data | Create a Sheet from CSV and read any sheet as header-keyed records |
| G-6 | Scriptable from the CLI | `gogoogle drive upload\|update\|comments`, `gogoogle docs suggestions`, `gogoogle sheets records` emit JSON |

### 2.2 Non-Goals

- Writing comments, replying, or resolving them.
- Accepting or rejecting suggestions.
- Sharing / permission management (callers share via Drive UI or existing
  tools).
- Drive search, folder sync, or storage abstractions.
- Rendering or converting Markdown (callers produce `.docx`/`.pptx`/`.csv`).
- Authentication changes beyond requesting the needed Drive scopes.

## 3. User Stories

**US-1 — Publish for review.** As an *application developer*, I want to
upload a generated `.docx` as a Google Doc into a given folder, so that
reviewers can comment and suggest in Google Docs.
Acceptance: `driveutil.Upload(ctx, UploadRequest{Path, Name, FolderID,
ConvertTo: Doc})` returns `File{ID, Name, MimeType, WebViewLink}`; a request
test asserts the target MIME type `application/vnd.google-apps.document`.

**US-2 — Publish a new version.** As an *application developer*, I want to
replace an existing reviewed file's contents with a new version, so that
reviewers keep the same link.
Acceptance: `driveutil.Update(ctx, fileID, path, convertTo)` keeps the file
ID; a request test asserts a media update on the existing ID.

**US-3 — Read reviewer comments.** As an *application or agent*, I want all
comments on a file with the text they refer to, so that each response can be
matched to the part of the document it concerns.
Acceptance: `driveutil.ListComments(ctx, fileID, opts)` returns all pages;
fields include `QuotedText`, `Resolved`, `Author`, `Replies`; deleted
comments are excluded unless requested.

**US-4 — Read proposed edits.** As an *application or agent*, I want a
document's suggested insertions and deletions separated from its original
text, so that I can propose the same edits to my source of truth.
Acceptance: `docsutil.ExtractSuggestionContent(doc)` over a fixture with
insertions and deletions in a paragraph and a table cell yields segments with
`Kind` ∈ {`original`, `insertion`, `deletion`} and suggestion IDs.

**US-5 — Review in a sheet.** As an *application developer*, I want to
create a review sheet from CSV and later read the edited sheet as records,
so that structured responses (agree / change / comment per row) come back as
data.
Acceptance: CSV upload converts to a Sheet; `sheetsutil.ReadRecords`
returns `[]map[string]string` keyed by the header row, skipping empty rows.

**US-6 — Do it from the shell.** As a *CLI user*, I want commands for the
above, so that scripts and AI agents can drive the round trip.
Acceptance: CLI tests cover argument validation and JSON output shape.

## 4. Functional Requirements

- **PRD-FR-1** `driveutil/v3` MUST provide a service constructed like the
  other gogoogle packages (from an `*http.Client`), and `GetFile` returning
  ID, name, MIME type, web link, modified time, and head revision ID.
- **PRD-FR-2** Upload MUST support conversion targets Doc, Sheet, Slides, and
  no conversion; an optional parent folder; and return the created file.
- **PRD-FR-3** Update MUST replace the media of an existing file ID, with the
  same conversion rules, and return the updated file.
- **PRD-FR-4** Comment listing MUST page through all results and return
  author, times, content, quoted anchor text, resolved flag, deleted flag, and
  replies (author, time, content, action such as `resolve` / `reopen`).
- **PRD-FR-5** Suggestion-aware extraction MUST walk body paragraphs and
  table cells of a document fetched with `SUGGESTIONS_INLINE` and return
  ordered segments with kind and suggestion IDs, plus a convenience view of
  each paragraph / cell as original and as-if-accepted text.
- **PRD-FR-6** Sheets: create from CSV (via Drive conversion) and
  `ReadRecords(spreadsheetID, sheetName)` keyed by the header row.
- **PRD-FR-7** CLI commands per G-6, JSON to stdout, errors to stderr,
  non-zero exit on failure.
- **PRD-FR-8** Existing APIs MUST remain backward compatible; new behavior is
  additive.

## 5. Success Metrics (required)

| Metric | Baseline | Target | Measurement Method |
|--------|----------|--------|--------------------|
| Review round-trip steps supported (upload-convert, update, comments, suggestions, sheet records) | 0 of 5 | 5 of 5 | Feature checklist in release notes; tests per step |
| Suggestion classification accuracy on fixtures | n/a (flattened) | 100 % of fixture segments correctly classified | Fixture tests with recorded `documents.get` responses |
| Comment fields returned | 0 | All fields in PRD-FR-4 present for 100 % of fixture comments; all pages fetched | Paged fake-server tests |
| Test coverage of new code | 0 % | ≥ 80 % | `go test -cover` on `driveutil`, `docsutil`, `sheetsutil` |
| Breaking changes to existing APIs | — | 0 | `apidiff` against v0.12.0 |

## 6. Out of Scope / Future

- Writing or resolving comments; accepting suggestions.
- Permission management and sharing.
- Native Google Slides generation improvements in `slidesutil`.
- MCP / agent tool surfaces (these live in wrapper projects that import
  gogoogle).
