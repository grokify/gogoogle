# Google Slides

The `slidesutil/v1` package provides utilities for creating and manipulating Google Slides presentations.

## Features

- Create presentations with title slides
- Add slides with various layouts
- Insert text boxes, tables, lines, and shapes
- Markdown-to-Slides conversion
- Batch update operations
- Extract content from existing presentations

## Quick Start

```go
import (
    "context"
    "google.golang.org/api/slides/v1"
    "google.golang.org/api/option"
    "github.com/grokify/gogoogle/slidesutil/v1"
)

// Create Slides service
service, err := slides.NewService(ctx, option.WithHTTPClient(httpClient))
if err != nil {
    log.Fatal(err)
}
```

## Creating Presentations

### New Presentation

```go
presentation := &slides.Presentation{
    Title: "My Presentation",
}

created, err := service.Presentations.Create(presentation).Do()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created: %s\n", created.PresentationId)
```

### Add Title Slide

```go
requests := []*slides.Request{
    slidesutil.CreateSlideRequest("slide1", slides.PredefinedLayoutTITLE),
    slidesutil.InsertTextRequest("slide1_title", "Welcome to My Presentation"),
    slidesutil.InsertTextRequest("slide1_subtitle", "An Introduction"),
}

_, err := service.Presentations.BatchUpdate(presentationID, &slides.BatchUpdatePresentationRequest{
    Requests: requests,
}).Do()
```

## Slide Layouts

Available predefined layouts:

| Layout | Description |
|--------|-------------|
| `BLANK` | Empty slide |
| `TITLE` | Title and subtitle |
| `TITLE_AND_BODY` | Title with body text |
| `TITLE_AND_TWO_COLUMNS` | Title with two columns |
| `TITLE_ONLY` | Title only |
| `SECTION_HEADER` | Section divider |
| `ONE_COLUMN_TEXT` | Single text column |
| `MAIN_POINT` | Large centered text |
| `BIG_NUMBER` | Featured number |

## Adding Content

### Text Box

```go
requests := []*slides.Request{
    slidesutil.CreateTextBoxRequest(
        "textbox1",
        "slideID",
        100, // x position (pt)
        100, // y position (pt)
        400, // width (pt)
        200, // height (pt)
        "Hello, World!",
    ),
}
```

### Table

```go
requests := []*slides.Request{
    slidesutil.CreateTableRequest(
        "table1",
        "slideID",
        3, // rows
        4, // columns
        100, // x
        150, // y
    ),
}
```

### Image

```go
requests := []*slides.Request{
    {
        CreateImage: &slides.CreateImageRequest{
            ObjectId: "image1",
            Url:      "https://example.com/image.png",
            ElementProperties: &slides.PageElementProperties{
                PageObjectId: "slideID",
                Size: &slides.Size{
                    Width:  &slides.Dimension{Magnitude: 300, Unit: "PT"},
                    Height: &slides.Dimension{Magnitude: 200, Unit: "PT"},
                },
                Transform: &slides.AffineTransform{
                    ScaleX:     1,
                    ScaleY:     1,
                    TranslateX: 100,
                    TranslateY: 100,
                    Unit:       "PT",
                },
            },
        },
    },
}
```

## Extracting Content

### Get Presentation Content

```go
content, err := slidesutil.ExtractPresentationContent(ctx, service, presentationID)
if err != nil {
    log.Fatal(err)
}

for _, slide := range content.Slides {
    fmt.Printf("Slide: %s\n", slide.Title)
    fmt.Printf("Content: %s\n", slide.TextContent)
    fmt.Printf("Notes: %s\n", slide.Notes)

    for _, img := range slide.Images {
        fmt.Printf("Image: %s\n", img.URL)
    }
}
```

### Extract Text Only

```go
text := slidesutil.ExtractTextContent(slide)
fmt.Println(text)
```

### Extract Images

```go
images := slidesutil.ExtractImages(slide)
for _, img := range images {
    fmt.Printf("Image URL: %s\n", img.URL)
    fmt.Printf("Content URL: %s\n", img.ContentURL)
}
```

## Markdown to Slides

Convert Markdown to presentation slides:

```go
markdown := `
# Introduction

Welcome to the presentation.

## Key Points

- Point one
- Point two
- Point three

## Conclusion

Thank you!
`

requests := slidesutil.MarkdownToSlideRequests(markdown)
```

## Styling

### Text Formatting

```go
requests := []*slides.Request{
    {
        UpdateTextStyle: &slides.UpdateTextStyleRequest{
            ObjectId: "textbox1",
            TextRange: &slides.Range{
                Type: "ALL",
            },
            Style: &slides.TextStyle{
                Bold:       true,
                FontSize:   &slides.Dimension{Magnitude: 24, Unit: "PT"},
                ForegroundColor: &slides.OptionalColor{
                    OpaqueColor: &slides.OpaqueColor{
                        RgbColor: &slides.RgbColor{
                            Red:   0.2,
                            Green: 0.4,
                            Blue:  0.8,
                        },
                    },
                },
            },
            Fields: "bold,fontSize,foregroundColor",
        },
    },
}
```

### Colors

```go
// RGB color
color := slidesutil.RGBColor(0.2, 0.4, 0.8) // Blue

// Hex to RGB
color := slidesutil.HexToRGB("#4285F4") // Google Blue
```

## CLI: Extract Content

Extract presentation content as JSON:

```bash
gogoogle slides content --presentation-id "1abc123..." --output content.json
```

## OAuth Scopes

| Scope | Description |
|-------|-------------|
| `presentations.readonly` | Read presentations |
| `presentations` | Read and write presentations |

## Best Practices

1. **Use batch updates** - Combine multiple changes into one API call
2. **Generate unique IDs** - Use UUIDs for element IDs
3. **Handle quotas** - Slides API has rate limits
4. **Cache presentations** - Avoid repeated fetches

## Next Steps

- [CLI Tools](../cli/index.md) - Command-line slides extraction
- [Sheets Integration](../sheets/index.md) - Data-driven presentations
