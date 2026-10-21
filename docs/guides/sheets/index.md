# Google Sheets

The `sheetsutil` package provides utilities for reading and writing Google Sheets data.

## Packages

| Package | Description |
|---------|-------------|
| `sheetsutil/v4` | Cell value parsing and URL utilities |
| `sheetsutil/v4/sheetsmap` | Map sheet data to Go types with validation |
| `sheetsutil/iwark` | Low-level operations using Iwark library |

## Cell Value Parsing

Parse Google Sheets API responses into typed Go structures:

```go
import sheetsutil "github.com/grokify/gogoogle/sheetsutil/v4"

// Get values from Sheets API
vr, _ := service.Spreadsheets.Values.Get(spreadsheetID, "Sheet1!A1:D10").Do()

// Parse to typed grid - each cell has type information
grid := sheetsutil.ParseValueRange(vr)
for _, row := range grid {
    for _, cell := range row {
        fmt.Printf("Type: %s, Value: %s\n", cell.Type, cell.FormattedValue)
    }
}

// Or extract as simple strings
formatted := sheetsutil.ExtractFormattedValues(vr)  // Display values
raw := sheetsutil.ExtractRawValues(vr)               // Underlying values
```

### Cell Types

| Type | Description |
|------|-------------|
| `CellTypeEmpty` | Empty cell |
| `CellTypeString` | Text value |
| `CellTypeNumber` | Numeric value |
| `CellTypeBool` | Boolean (TRUE/FALSE) |
| `CellTypeDate` | Date value |
| `CellTypeTime` | Time value |
| `CellTypeDateTime` | Date and time |
| `CellTypeDuration` | Duration |
| `CellTypeError` | Error value |

## URL Utilities

Parse spreadsheet IDs from URLs and build URLs:

```go
import sheetsutil "github.com/grokify/gogoogle/sheetsutil/v4"

// Extract spreadsheet ID from URL or plain ID
id, err := sheetsutil.ParseSpreadsheetURL(
    "https://docs.google.com/spreadsheets/d/abc123/edit#gid=0",
)
// id = "abc123"

// Extract full info including sheet GID and range
info, err := sheetsutil.ParseSpreadsheetURLFull(
    "https://docs.google.com/spreadsheets/d/abc123/edit#gid=456&range=A1:D10",
)
// info.SpreadsheetID = "abc123"
// info.SheetGID = 456
// info.Range = "A1:D10"

// Build URLs from ID
url := sheetsutil.BuildSpreadsheetURL("abc123")
// "https://docs.google.com/spreadsheets/d/abc123/edit"

url := sheetsutil.BuildSpreadsheetURLWithSheet("abc123", 456)
// "https://docs.google.com/spreadsheets/d/abc123/edit#gid=456"

// Check if string is a Sheets URL
if sheetsutil.IsSpreadsheetURL(input) {
    // Handle URL
}
```

## Quick Start

```go
import (
    "context"
    "google.golang.org/api/sheets/v4"
    "google.golang.org/api/option"
)

// Create Sheets service
service, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
if err != nil {
    log.Fatal(err)
}

// Read spreadsheet
spreadsheetID := "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"
readRange := "Sheet1!A1:D10"

resp, err := service.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
if err != nil {
    log.Fatal(err)
}

for _, row := range resp.Values {
    fmt.Println(row)
}
```

## sheetsmap

Map spreadsheet data to typed Go structures:

```go
import "github.com/grokify/gogoogle/sheetsutil/v4/sheetsmap"

// Define your data structure
type Contact struct {
    Name  string
    Email string
    Phone string
}

// Read and map data
data, err := sheetsmap.ReadSpreadsheet(ctx, service, spreadsheetID, "Contacts")
if err != nil {
    log.Fatal(err)
}

// Access rows as maps
for _, row := range data.Rows {
    name := row["Name"].(string)
    email := row["Email"].(string)
    fmt.Printf("%s: %s\n", name, email)
}
```

### Column Mapping

```go
// Define column mappings
columns := sheetsmap.Columns{
    {Name: "Name", Required: true},
    {Name: "Email", Required: true, Validator: sheetsmap.ValidateEmail},
    {Name: "Status", Enum: []string{"active", "inactive", "pending"}},
}

// Read with validation
data, err := sheetsmap.ReadSpreadsheetWithColumns(ctx, service, spreadsheetID, "Sheet1", columns)
```

### Enum Validation

```go
columns := sheetsmap.Columns{
    {
        Name: "Priority",
        Enum: []string{"low", "medium", "high", "critical"},
    },
    {
        Name: "Status",
        Enum: []string{"open", "in_progress", "resolved", "closed"},
    },
}
```

## Writing Data

```go
// Prepare values
values := [][]interface{}{
    {"Name", "Email", "Phone"},
    {"John Doe", "john@example.com", "555-0100"},
    {"Jane Smith", "jane@example.com", "555-0101"},
}

// Write to sheet
_, err := service.Spreadsheets.Values.Update(
    spreadsheetID,
    "Sheet1!A1",
    &sheets.ValueRange{Values: values},
).ValueInputOption("RAW").Do()
```

### Append Rows

```go
newRows := [][]interface{}{
    {"New User", "new@example.com", "555-0199"},
}

_, err := service.Spreadsheets.Values.Append(
    spreadsheetID,
    "Sheet1!A1",
    &sheets.ValueRange{Values: newRows},
).ValueInputOption("USER_ENTERED").Do()
```

## iwark Integration

Low-level operations using [Iwark spreadsheet](https://github.com/Iwark/spreadsheet):

```go
import "github.com/grokify/gogoogle/sheetsutil/iwark"

// Fetch spreadsheet
spreadsheet, err := iwark.FetchSpreadsheet(ctx, httpClient, spreadsheetID)
if err != nil {
    log.Fatal(err)
}

// Access sheets
for _, sheet := range spreadsheet.Sheets {
    fmt.Printf("Sheet: %s\n", sheet.Properties.Title)
}
```

## OAuth Scopes

| Scope | Description |
|-------|-------------|
| `spreadsheets.readonly` | Read spreadsheets |
| `spreadsheets` | Read and write spreadsheets |

```go
scopes := []string{
    "https://www.googleapis.com/auth/spreadsheets.readonly",
}
```

## Common Operations

### Get Spreadsheet Metadata

```go
spreadsheet, err := service.Spreadsheets.Get(spreadsheetID).Do()
if err != nil {
    log.Fatal(err)
}

for _, sheet := range spreadsheet.Sheets {
    fmt.Printf("Sheet: %s (ID: %d)\n",
        sheet.Properties.Title,
        sheet.Properties.SheetId)
}
```

### Create New Spreadsheet

```go
spreadsheet := &sheets.Spreadsheet{
    Properties: &sheets.SpreadsheetProperties{
        Title: "My New Spreadsheet",
    },
}

created, err := service.Spreadsheets.Create(spreadsheet).Do()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created: %s\n", created.SpreadsheetUrl)
```

### Add New Sheet

```go
requests := []*sheets.Request{
    {
        AddSheet: &sheets.AddSheetRequest{
            Properties: &sheets.SheetProperties{
                Title: "New Tab",
            },
        },
    },
}

_, err := service.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
    Requests: requests,
}).Do()
```

## Best Practices

1. **Use ValueInputOption wisely**
   - `RAW`: Values entered as-is
   - `USER_ENTERED`: Parse as if typed by user (formulas work)

2. **Batch operations** - Combine multiple updates

3. **Handle rate limits** - Google Sheets API has quotas

4. **Validate data** - Use sheetsmap validators

## Next Steps

- [Gmail Mail Merge](../gmail/mail-merge.md) - Use Sheets data for email campaigns
- [CLI Tools](../cli/index.md) - Command-line utilities
