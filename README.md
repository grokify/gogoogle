# Go Google

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/grokify/gogoogle/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/grokify/gogoogle/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/grokify/gogoogle/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/grokify/gogoogle/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/grokify/gogoogle/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/grokify/gogoogle/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/gogoogle
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/gogoogle
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/gogoogle
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/gogoogle
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fgogoogle
 [loc-svg]: https://tokei.rs/b1/github/grokify/gogoogle
 [repo-url]: https://github.com/grokify/gogoogle
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/gogoogle/blob/master/LICENSE

`gogoogle` is a set of generic, higher-level libraries for interacting with Google APIs using Go. It is built on the official [Google API Go Client](https://github.com/google/google-api-go-client) and [Google Cloud Go](https://github.com/googleapis/google-cloud-go) libraries.

## Installation

```bash
go get github.com/grokify/gogoogle
```

## Packages

### Gmail (`gmailutil/v1`)

Comprehensive Gmail API helper library for message operations, label management, and mail merge functionality.

- List and retrieve messages with filtering options (by sender, category, labels)
- Batch delete messages
- Send emails with the `mailutil.MessageWriter` interface
- Mail merge with Mustache templates and Google Sheets data sources
- Label management

### Google Sheets (`sheetsutil`)

Utilities for reading and writing Google Sheets data with typed structures.

- **sheetsutil/v4/sheetsmap** - Maps sheet data to Go types with enum validation and column management
- **sheetsutil/iwark** - Low-level spreadsheet operations using the [Iwark spreadsheet](https://github.com/Iwark/spreadsheet) library

### Google Slides (`slidesutil/v1`)

Comprehensive Google Slides manipulation library.

- Create presentations with title slides
- Add slides with various layouts (title/body, main point, image sidebars)
- Create and style text boxes, tables, lines, and shapes
- Markdown-to-Slides conversion with bulleted lists, text styling, and links
- Batch update operations
- Color utilities for RGB/hex conversion

### Google Maps (`mapsutil/staticmap`)

Generate Google Static Maps with customizable markers.

- Configure map center, zoom, and dimensions
- Add styled markers with colors and labels
- Preset regions (USA, Europe, World)
- Download maps as PNG files

### Speech-to-Text (`speechtotext`)

Wrapper for Google's Speech-to-Text API.

- Transcribe audio from byte data or files
- Confidence threshold filtering

### Text-to-Speech (`texttospeech/v1beta1`)

Wrapper for Google's Text-to-Speech API.

- Synthesize speech from text
- Multiple voice options (WaveNet neural, standard)
- Multiple audio formats (MP3, LINEAR16, OGG_OPUS)

### Data Loss Prevention (`dlp/v2`)

Helper functions for Google's DLP API to detect sensitive information.

- Create content items for DLP inspection
- Pre-defined info types (credit cards, SSNs, person names, US states)

### BigQuery (`bigqueryutil`)

Utilities for streaming data uploads to BigQuery.

- Automatic chunking for large uploads (max 10,000 items per operation)
- Error handling for batch insert operations

### Google Docs (`docsutil`)

URL parsing and construction for Google Docs and Sheets.

- Parse document IDs from URLs
- Build document URLs from IDs

### Authentication (`auth`)

Simplified OAuth2 setup for Google APIs.

- Token management with file-based storage
- Token refresh support

## Related Libraries

- OAuth 2.0 utilities via [`goauth/google`](https://github.com/grokify/goauth/tree/master/google)
- OAuth 2.0 demo app via [`beegoutil`](https://github.com/grokify/beegoutil)
