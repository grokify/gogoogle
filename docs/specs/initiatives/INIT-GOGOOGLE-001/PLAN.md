# PLAN — Document review round trip — Drive upload and conversion, comments, and suggestion-aware Docs content

**Initiative:** `INIT-GOGOOGLE-001`
**Repository:** `github.com/grokify/gogoogle`
**Version:** 0.1
**Status:** Draft — pending tech-lead review

Non-normative sequencing of [PRD.md](PRD.md) and [TRD.md](TRD.md); execution
state lives in [ROADMAP.md](ROADMAP.md).

## 1. Approach

The Drive package comes first. Upload and comments are gaps no existing
package can fill, and both conversion and CSV-to-Sheet creation depend on it.

Suggestion-aware extraction is independent of Drive and runs in parallel.

The CLI and docs come last, wrapping finished library functions. Each phase
ends in a usable, releasable state.

## 2. Affected Components

| Component | Change |
|-----------|--------|
| `driveutil/v3` | New |
| `docsutil/v1` | Extended (new function and types; existing untouched) |
| `sheetsutil/v4` | Extended (`ReadRecords`, `CreateFromCSV`) |
| `cmd/gogoogle` | New `drive`, `docs`, `sheets` subcommands |
| `docs/guides`, `docs/specs`, `mkdocs.yml`, `CHANGELOG.json` | Updated |

## 3. Work Breakdown

| Step | Approach | Satisfies | Verification | Depends on |
|------|----------|-----------|--------------|------------|
| P-01 | `driveutil/v3` service and `GetFile` | PRD-FR-1; TRD-001 | Fake-server test | — |
| P-02 | `Upload` / `Update` with conversion | PRD-FR-2, 3; TRD-002, 003 | Request-capture tests | P-01 |
| P-03 | `ListComments` with paging and replies | PRD-FR-4; TRD-004 | Paged fake-server test | P-01 |
| P-04 | `docsutil.ExtractSuggestionContent` | PRD-FR-5; TRD-005 | Recorded fixtures; INV-4 property test | — |
| P-05 | `sheetsutil.ReadRecords` and `CreateFromCSV` | PRD-FR-6; TRD-006 | Fake-server test | P-02 |
| P-06 | CLI commands | PRD-FR-7; TRD-007 | CLI tests | P-02, P-03, P-04, P-05 |
| P-07 | Guides, specs index, `apidiff` check, CHANGELOG, release | PRD-FR-8; TRD-008 | CI green, then tag | P-06 |

## 4. Sequencing and Dependencies

```
P-01 → {P-02, P-03};  P-02 → P-05;  P-04 (parallel);  {P-02..P-05} → P-06 → P-07
```

**Critical path:** P-01 → P-02 → P-05 → P-06 → P-07.

**Parallel:** P-03 and P-04 run alongside P-02.

## 5. Data and Migration

Not applicable — additive library changes.

## 6. Testing Strategy

Per TRD §12. Coverage target ≥ 80 % on new code.

## 7. Rollout and Rollback

- **Release.** Ship as one minor release (v0.13.0) once Phase 3 completes.
- **Optional pre-release.** Cut one after Phase 1 for downstream
  integration.
- **Tagging.** Tag only after CI is green.
- **Rollback.** Pin v0.12.0.

## 8. Risks

| ID | Risk | Likelihood / Impact | Mitigation |
|----|------|---------------------|------------|
| PR-1 | Recorded fixtures drift from the live API shape | Low / Medium | Optional integration test refreshes fixtures |
| PR-2 | Scope confusion (`drive.file` vs `drive.readonly`) | Medium / Low | Guide documents scope per function |
