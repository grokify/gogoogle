# ROADMAP — Document review round trip — Drive upload and conversion, comments, and suggestion-aware Docs content

**Initiative:** `INIT-GOGOOGLE-001`
**Repository:** `github.com/grokify/gogoogle`
**Status:** Planned — 0 of 4 phases completed

> RMI IDs are stable and permanent. Commits implementing an item carry the trailer `Refs: RMI-GOGOOGLE-<NNN>`. Phase status is derived from member RMIs. RMI `NNN` corresponds to PLAN step `P-NN`.

## Phase 1 — Drive Package

**Theme:** Upload with conversion, update in place, and comment listing through a new `driveutil/v3` package
**Status:** Planned — 0 of 3 items completed

- [ ] `RMI-GOGOOGLE-001` `driveutil/v3` service construction and GetFile
  - Type: capability · Status: planned
  - Satisfies: PRD-FR-1, TRD-001
  - Acceptance: fake-server test returns ID, name, MIME type, web link, modified time, head revision
- [ ] `RMI-GOGOOGLE-002` Upload with conversion (Doc / Sheet / Slides) and update in place
  - Type: capability · Status: planned
  - Depends on: `RMI-GOGOOGLE-001`
  - Satisfies: PRD-FR-2, PRD-FR-3, TRD-002, TRD-003
  - Acceptance: request-capture tests per conversion; update keeps file ID
- [ ] `RMI-GOGOOGLE-003` Comment listing with quoted anchor text, resolved state, and replies
  - Type: capability · Status: planned
  - Depends on: `RMI-GOGOOGLE-001`
  - Satisfies: PRD-FR-4, TRD-004
  - Acceptance: multi-page fake-server test returns all comments with all fields

## Phase 2 — Review Content

**Theme:** Suggestion-aware Docs extraction and Sheets records for structured review responses
**Status:** Planned — 0 of 2 items completed

- [ ] `RMI-GOGOOGLE-004` Suggestion-aware Docs content extraction (original / insertion / deletion segments with original and accepted views)
  - Type: capability · Status: planned
  - Satisfies: PRD-FR-5, TRD-005
  - Acceptance: recorded fixtures with paragraph and table-cell suggestions classify 100 % of segments; INV-4 holds
- [ ] `RMI-GOGOOGLE-005` Sheets ReadRecords and CreateFromCSV
  - Type: capability · Status: planned
  - Depends on: `RMI-GOGOOGLE-002`
  - Satisfies: PRD-FR-6, TRD-006
  - Acceptance: header-keyed records; CSV converts to a Sheet

## Phase 3 — CLI, Docs, and Release

**Theme:** Scriptable commands for the round trip, guides, compatibility check, and release
**Status:** Planned — 0 of 2 items completed

- [ ] `RMI-GOGOOGLE-006` CLI: drive upload / update / comments, docs suggestions, sheets records
  - Type: capability · Status: planned
  - Depends on: `RMI-GOGOOGLE-002`, `RMI-GOGOOGLE-003`, `RMI-GOGOOGLE-004`, `RMI-GOGOOGLE-005`
  - Satisfies: PRD-FR-7, TRD-007
  - Acceptance: CLI tests cover usage errors and JSON output
- [ ] `RMI-GOGOOGLE-007` Guides, specs index, apidiff compatibility check, CHANGELOG, and release
  - Type: release · Status: planned
  - Depends on: `RMI-GOGOOGLE-006`
  - Satisfies: PRD-FR-8, TRD-008
  - Acceptance: no breaking API changes; CI green before tag

## Phase 4 — Confined Uploads

**Theme:** Send content from a reader the caller already opened and checked, so a path is not reopened after it is verified
**Status:** Planned — 0 of 1 items completed

- [ ] `RMI-GOGOOGLE-008` Upload and update from an io.Reader
  - Type: capability · Status: planned
  - Depends on: `RMI-GOGOOGLE-002`
  - Acceptance: upload and update send a reader's content without opening a path; a reader error is returned wrapped; existing path-based calls behave as before
