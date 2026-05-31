# SheetsMap

SheetsMap treats a Google Sheet as a key-value store with typed columns and enum validation.

## Overview

SheetsMap provides:

- **Key-value storage** - First column is the key, second is display name, rest are data columns
- **Column schema** - Define columns with names, aliases, abbreviations, and enum values
- **Enum validation** - Validate and canonicalize values against allowed enums
- **CRUD operations** - Get, create, update items with automatic synchronization

## Sheet Structure

The first row defines the schema. Each cell follows this format:

```
column_name|alias1|alias2 - abbreviation - enum1|alt1, enum2|alt2 - info_text|info_url
```

**Example header row:**

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| `email` | `name` | `tshirt size\|t-shirt size - Size - XS, S, M, L, XL` |

## Usage

### Initialize

```go
import (
    "github.com/grokify/gogoogle/sheetsutil/v4/sheetsmap"
    "google.golang.org/api/sheets/v4"
)

// Using sheet index
sm, err := sheetsmap.NewSheetsMapIndex(httpClient, spreadsheetID, 0)

// Using sheet title
sm, err := sheetsmap.NewSheetsMapTitle(httpClient, spreadsheetID, "Sheet1")
```

### Read Data

```go
// Read columns and items
err := sm.FullRead()

// Access items by key
item, err := sm.GetItem("user@example.com")
fmt.Println(item.Display)      // Display name
fmt.Println(item.Data["size"]) // Column value
```

### Write Data

```go
// Get or create an item
item, err := sm.GetOrCreateItem("user@example.com")

// Update a field with enum validation
_, err = sm.UpdateItem(item, "tshirt size", "M", true)

// Set item display name
err = sm.SetItemKeyDisplay("user@example.com", "John Doe")
```

### Column Schema

```go
// Parse a column definition
col, err := sheetsmap.ParseColumn("size|t-shirt size - S - XS, S, M, L, XL")

// col.Name = "size"
// col.NameAliases = ["t-shirt size"]
// col.Abbreviation = "S"
// col.Enums = [{Canonical: "XS"}, {Canonical: "S"}, ...]

// Validate and canonicalize a value
canonical, err := col.ValueToCanonical("extra small") // Returns "XS"
```

## Column Definition Format

```
name|alias1|alias2 - abbreviation - enum1|alt1|alt2, enum2|alt1 - label|url ~ label2|url2
```

| Part | Description | Example |
|------|-------------|---------|
| Name | Primary column name | `tshirt size` |
| Aliases | Alternative names (pipe-separated) | `t-shirt size\|shirt size` |
| Abbreviation | Short form | `Size` |
| Enums | Allowed values with aliases | `XS\|extra small, S\|small` |
| Info URLs | Reference links | `Size Chart\|http://example.com` |

## Authentication

Use any `*http.Client` with Google OAuth credentials:

```go
import "github.com/grokify/goauth/google"

client, err := google.NewClientFromJWTJSON(
    ctx,
    []byte(serviceAccountJSON),
    sheets.DriveScope,
    sheets.SpreadsheetsScope,
)

sm, err := sheetsmap.NewSheetsMapIndex(client, spreadsheetID, 0)
```
