# Speech-to-Text

The `speechtotext` package provides a wrapper for Google's Speech-to-Text API.

## Features

- Transcribe audio from bytes or files
- Multiple language support
- Confidence threshold filtering
- Word-level timestamps

## Quick Start

```go
import (
    "context"
    "os"
    "github.com/grokify/gogoogle/speechtotext"
)

// Read audio file
audioData, _ := os.ReadFile("audio.wav")

// Transcribe
result, err := speechtotext.Transcribe(ctx, httpClient, speechtotext.TranscribeRequest{
    Audio:        audioData,
    LanguageCode: "en-US",
    Encoding:     "LINEAR16",
    SampleRate:   16000,
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(result.Transcript)
```

## TranscribeRequest

| Field | Type | Description |
|-------|------|-------------|
| `Audio` | `[]byte` | Audio data |
| `LanguageCode` | `string` | BCP-47 language code |
| `Encoding` | `string` | Audio encoding |
| `SampleRate` | `int` | Sample rate in Hz |
| `MinConfidence` | `float64` | Minimum confidence (0-1) |

## Audio Encodings

| Encoding | Description |
|----------|-------------|
| `LINEAR16` | Uncompressed 16-bit PCM |
| `FLAC` | FLAC encoded |
| `MULAW` | μ-law encoded |
| `AMR` | AMR (Adaptive Multi-Rate) |
| `AMR_WB` | AMR Wideband |
| `OGG_OPUS` | Ogg Opus |
| `MP3` | MP3 encoded |

## Language Codes

Common language codes:

| Language | Code |
|----------|------|
| English (US) | `en-US` |
| English (UK) | `en-GB` |
| Spanish | `es-ES` |
| French | `fr-FR` |
| German | `de-DE` |
| Japanese | `ja-JP` |
| Chinese (Mandarin) | `zh-CN` |
| Portuguese (Brazil) | `pt-BR` |

## Transcription Result

```go
type TranscribeResult struct {
    Transcript  string    // Full transcript
    Confidence  float64   // Overall confidence (0-1)
    Words       []Word    // Word-level results
}

type Word struct {
    Word       string
    StartTime  time.Duration
    EndTime    time.Duration
    Confidence float64
}
```

## Confidence Filtering

```go
result, err := speechtotext.Transcribe(ctx, httpClient, speechtotext.TranscribeRequest{
    Audio:         audioData,
    LanguageCode:  "en-US",
    MinConfidence: 0.8, // Only include high-confidence results
})
```

## File Transcription

```go
result, err := speechtotext.TranscribeFile(ctx, httpClient, "audio.wav", speechtotext.TranscribeRequest{
    LanguageCode: "en-US",
})
```

## Long Audio

For audio longer than 1 minute, use async recognition:

```go
operation, err := speechtotext.TranscribeAsync(ctx, httpClient, speechtotext.AsyncTranscribeRequest{
    AudioURI:     "gs://bucket/audio.wav", // GCS URI
    LanguageCode: "en-US",
})

// Poll for completion
result, err := speechtotext.WaitForResult(ctx, httpClient, operation)
```

## OAuth Scope

```go
scope := "https://www.googleapis.com/auth/cloud-platform"
```

## Enable API

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Enable **Cloud Speech-to-Text API**
3. Create service account or OAuth credentials

## Best Practices

1. **Use appropriate encoding** - LINEAR16 for highest accuracy
2. **Set correct sample rate** - Must match audio
3. **Specify language** - Don't rely on auto-detection
4. **Use GCS for long audio** - Required for files > 1 minute

## Error Handling

```go
result, err := speechtotext.Transcribe(ctx, httpClient, request)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "INVALID_ARGUMENT"):
        log.Println("Invalid audio format or parameters")
    case strings.Contains(err.Error(), "RESOURCE_EXHAUSTED"):
        log.Println("Quota exceeded")
    default:
        log.Printf("Transcription error: %v", err)
    }
}
```

## Next Steps

- [Text-to-Speech](tts.md) - Convert text to audio
- [Gmail](../gmail/index.md) - Send transcripts via email
