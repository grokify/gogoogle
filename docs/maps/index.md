# Google Maps

The `mapsutil/staticmap` package generates Google Static Maps images.

## Overview

Create map images with:

- Custom center, zoom, and dimensions
- Styled markers with colors and labels
- Preset regions (USA, Europe, World)
- Download as PNG files

## Quick Start

```go
import "github.com/grokify/gogoogle/mapsutil/staticmap"

// Create a map centered on San Francisco
mapConfig := staticmap.MapConfig{
    Center:  "San Francisco, CA",
    Zoom:    12,
    Width:   600,
    Height:  400,
    MapType: "roadmap",
    APIKey:  os.Getenv("GOOGLE_MAPS_API_KEY"),
}

url := staticmap.BuildURL(mapConfig)
fmt.Println(url)
```

## MapConfig Options

| Field | Type | Description |
|-------|------|-------------|
| `Center` | `string` | Map center (address or lat,lng) |
| `Zoom` | `int` | Zoom level (1-21) |
| `Width` | `int` | Image width in pixels |
| `Height` | `int` | Image height in pixels |
| `MapType` | `string` | Map type |
| `Markers` | `[]Marker` | Map markers |
| `APIKey` | `string` | Google Maps API key |

### Map Types

| Type | Description |
|------|-------------|
| `roadmap` | Standard road map |
| `satellite` | Satellite imagery |
| `terrain` | Terrain with roads |
| `hybrid` | Satellite with roads |

## Adding Markers

### Basic Marker

```go
mapConfig := staticmap.MapConfig{
    Center: "New York, NY",
    Zoom:   13,
    Width:  600,
    Height: 400,
    Markers: []staticmap.Marker{
        {
            Location: "Times Square, New York",
        },
    },
}
```

### Styled Markers

```go
markers := []staticmap.Marker{
    {
        Location: "40.758,-73.985",
        Color:    "red",
        Label:    "A",
    },
    {
        Location: "40.748,-73.986",
        Color:    "blue",
        Label:    "B",
    },
    {
        Location: "40.752,-73.977",
        Color:    "green",
        Label:    "C",
    },
}
```

### Marker Options

| Field | Type | Description |
|-------|------|-------------|
| `Location` | `string` | Address or lat,lng |
| `Color` | `string` | Marker color |
| `Label` | `string` | Single character label |
| `Size` | `string` | Marker size |

#### Colors

Standard colors: `red`, `blue`, `green`, `yellow`, `purple`, `orange`, `gray`, `white`, `black`

Custom hex: `0xFF5733`

#### Sizes

- `tiny` - Smallest
- `small` - Small
- `mid` - Medium (default)

## Preset Regions

```go
// USA map
usaConfig := staticmap.USAMapConfig()
usaConfig.APIKey = apiKey
usaConfig.Markers = markers

// Europe map
europeConfig := staticmap.EuropeMapConfig()

// World map
worldConfig := staticmap.WorldMapConfig()
```

## Download Map

```go
// Get map image as bytes
imageData, err := staticmap.Download(mapConfig)
if err != nil {
    log.Fatal(err)
}

// Save to file
err = os.WriteFile("map.png", imageData, 0644)
```

## URL Building

Generate URL without downloading:

```go
url := staticmap.BuildURL(mapConfig)
fmt.Println(url)
// https://maps.googleapis.com/maps/api/staticmap?center=...
```

## Multiple Markers Example

```go
import "github.com/grokify/gogoogle/mapsutil/staticmap"

func createOfficeMap() error {
    offices := []staticmap.Marker{
        {Location: "San Francisco, CA", Color: "red", Label: "H"},    // HQ
        {Location: "New York, NY", Color: "blue", Label: "N"},
        {Location: "Austin, TX", Color: "blue", Label: "A"},
        {Location: "Seattle, WA", Color: "blue", Label: "S"},
    }

    mapConfig := staticmap.MapConfig{
        Center:   "United States",
        Zoom:     4,
        Width:    800,
        Height:   500,
        MapType:  "roadmap",
        Markers:  offices,
        APIKey:   os.Getenv("GOOGLE_MAPS_API_KEY"),
    }

    imageData, err := staticmap.Download(mapConfig)
    if err != nil {
        return err
    }

    return os.WriteFile("offices.png", imageData, 0644)
}
```

## Zoom Levels

| Zoom | Coverage |
|------|----------|
| 1 | World |
| 5 | Continent |
| 10 | City |
| 15 | Streets |
| 20 | Buildings |

## API Key

Get a Google Maps API key:

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Enable **Maps Static API**
3. Create credentials (API key)
4. Restrict key to Static Maps API

```bash
export GOOGLE_MAPS_API_KEY="your-api-key"
```

## Limits

- **Free tier**: First $200/month free
- **Image size**: Max 640x640 (free) or 2048x2048 (premium)
- **Markers**: Max ~100 per request (URL length limit)

## Best Practices

1. **Secure API key** - Restrict to specific APIs and IPs
2. **Cache images** - Don't regenerate static maps unnecessarily
3. **Optimize size** - Use appropriate dimensions for your use case
4. **Handle errors** - API may return errors for invalid locations

## Next Steps

- [Gmail](../gmail/index.md) - Email maps as attachments
- [Sheets](../sheets/index.md) - Generate maps from location data
