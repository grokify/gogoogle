# Changelog

See [CHANGELOG.md](https://github.com/grokify/gogoogle/blob/main/CHANGELOG.md) for the full changelog.

## Releases

| Version | Date | Highlights |
|---------|------|------------|
| v0.9.0 | 2026-03-22 | SendSimple for simplified email sending |
| v0.8.0 | 2026-02-14 | Unified gogoogle CLI |
| v0.7.0 | 2026-02-09 | Slides content extraction |
| v0.6.0 | 2025-01-21 | Gmail mail merge with templates |
| v0.5.0 | 2025-01-08 | Linting and documentation |
| v0.4.0 | 2024-10-06 | Module renamed to gogoogle |

## Latest Release: v0.9.0

### Added

- `gmailutil/v1`: `SendSimple()` method for sending emails with minimal configuration
- `gmailutil/v1`: `SendSimpleOpts{}` struct with To, Subject, BodyText, BodyHTML, and ReplyTo fields

### Fixed

- `gmailutil/v1`: fix typo in `Send()` comment

## v0.8.0

### Added

- `cmd/gogoogle`: unified Cobra CLI for Google APIs with slides and gmail subcommands
- `cmd/gogoogle slides content`: extract text, images, and notes from Google Slides presentations
- `cmd/gogoogle gmail merge`: send templated emails via mail merge with Google Sheets data
- `cmd/gogoogle gmail send-markdown`: send emails with markdown body converted to HTML

## v0.7.0

### Added

- `slidesutil/v1`: `ExtractPresentationContent()` for extracting all slides' text and images
- `slidesutil/v1`: Slide text extraction functions
- `slidesutil/v1`: Image extraction with `ImageInfo` struct

## v0.6.0

### Added

- `gmailutil/v1/mailmerge`: `MailMerge{}` and `MailMergeOpts{}` for template-based email campaigns
- `gmailutil/v1/mailmerge`: `ExecMailMergeCLI()` for CLI usage

## v0.5.0

### Documentation

- README updates and lint status badge

## v0.4.0

### Changed

- **BREAKING**: Module renamed from `github.com/grokify/googleutil` to `github.com/grokify/gogoogle`

## Versioning

GoGoogle follows [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking API changes
- **MINOR**: New features, backward compatible
- **PATCH**: Bug fixes, backward compatible
