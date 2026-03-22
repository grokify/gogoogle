# Text-to-Speech

The `texttospeech/v1beta1` package provides a wrapper for Google's Text-to-Speech API.

## Features

- Synthesize speech from text
- Multiple voice options (WaveNet, Standard)
- Multiple audio formats
- SSML support

## Quick Start

```go
import (
    "context"
    "os"
    "github.com/grokify/gogoogle/texttospeech/v1beta1"
)

// Synthesize speech
audio, err := texttospeech.Synthesize(ctx, httpClient, texttospeech.SynthesizeRequest{
    Text:         "Hello, world!",
    LanguageCode: "en-US",
    VoiceName:    "en-US-Wavenet-D",
    AudioFormat:  "MP3",
})
if err != nil {
    log.Fatal(err)
}

// Save to file
err = os.WriteFile("output.mp3", audio, 0644)
```

## SynthesizeRequest

| Field | Type | Description |
|-------|------|-------------|
| `Text` | `string` | Text to synthesize |
| `SSML` | `string` | SSML markup (alternative to Text) |
| `LanguageCode` | `string` | BCP-47 language code |
| `VoiceName` | `string` | Voice name |
| `SsmlGender` | `string` | Voice gender |
| `AudioFormat` | `string` | Output audio format |
| `SpeakingRate` | `float64` | Speaking rate (0.25-4.0) |
| `Pitch` | `float64` | Pitch (-20 to 20) |

## Voices

### WaveNet Voices (Neural)

High-quality, natural-sounding voices:

| Voice | Gender | Description |
|-------|--------|-------------|
| `en-US-Wavenet-A` | Male | US English |
| `en-US-Wavenet-B` | Male | US English |
| `en-US-Wavenet-C` | Female | US English |
| `en-US-Wavenet-D` | Male | US English |
| `en-US-Wavenet-E` | Female | US English |
| `en-US-Wavenet-F` | Female | US English |
| `en-GB-Wavenet-A` | Female | UK English |
| `en-GB-Wavenet-B` | Male | UK English |

### Standard Voices

Lower latency, lower cost:

| Voice | Gender | Description |
|-------|--------|-------------|
| `en-US-Standard-A` | Male | US English |
| `en-US-Standard-B` | Male | US English |
| `en-US-Standard-C` | Female | US English |
| `en-US-Standard-D` | Male | US English |

### List Available Voices

```go
voices, err := texttospeech.ListVoices(ctx, httpClient, "en-US")
for _, voice := range voices {
    fmt.Printf("%s (%s)\n", voice.Name, voice.SsmlGender)
}
```

## Audio Formats

| Format | Description |
|--------|-------------|
| `MP3` | MP3 audio |
| `LINEAR16` | Uncompressed WAV |
| `OGG_OPUS` | Ogg Opus |

## SSML Support

Use SSML for advanced speech control:

```go
ssml := `<speak>
    <say-as interpret-as="date" format="mdy">12/25/2024</say-as>
    <break time="500ms"/>
    <emphasis level="strong">Merry Christmas!</emphasis>
</speak>`

audio, err := texttospeech.Synthesize(ctx, httpClient, texttospeech.SynthesizeRequest{
    SSML:         ssml,
    LanguageCode: "en-US",
    VoiceName:    "en-US-Wavenet-D",
    AudioFormat:  "MP3",
})
```

### SSML Elements

| Element | Description |
|---------|-------------|
| `<break>` | Insert pause |
| `<emphasis>` | Emphasize text |
| `<say-as>` | Interpret as date, time, etc. |
| `<prosody>` | Control pitch, rate, volume |
| `<sub>` | Pronunciation substitution |

## Speaking Rate and Pitch

```go
audio, err := texttospeech.Synthesize(ctx, httpClient, texttospeech.SynthesizeRequest{
    Text:         "This is spoken slowly with a lower pitch.",
    LanguageCode: "en-US",
    VoiceName:    "en-US-Wavenet-D",
    AudioFormat:  "MP3",
    SpeakingRate: 0.8,  // Slower (default 1.0)
    Pitch:        -2.0, // Lower pitch (default 0)
})
```

## Languages

| Language | Code | Voices |
|----------|------|--------|
| English (US) | `en-US` | 30+ |
| English (UK) | `en-GB` | 10+ |
| Spanish | `es-ES` | 10+ |
| French | `fr-FR` | 10+ |
| German | `de-DE` | 10+ |
| Japanese | `ja-JP` | 8+ |
| Chinese | `cmn-CN` | 8+ |

## OAuth Scope

```go
scope := "https://www.googleapis.com/auth/cloud-platform"
```

## Enable API

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Enable **Cloud Text-to-Speech API**
3. Create credentials

## Pricing

| Voice Type | Price per 1M chars |
|------------|-------------------|
| Standard | $4.00 |
| WaveNet | $16.00 |
| Neural2 | $16.00 |

First 1M standard characters free per month.

## Best Practices

1. **Use WaveNet for quality** - More natural sounding
2. **Cache audio** - Don't regenerate unchanged text
3. **Use SSML** - Better control over pronunciation
4. **Choose appropriate format** - MP3 for web, LINEAR16 for processing

## Error Handling

```go
audio, err := texttospeech.Synthesize(ctx, httpClient, request)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "INVALID_ARGUMENT"):
        log.Println("Invalid voice or language")
    case strings.Contains(err.Error(), "RESOURCE_EXHAUSTED"):
        log.Println("Quota exceeded")
    default:
        log.Printf("TTS error: %v", err)
    }
}
```

## Next Steps

- [Speech-to-Text](stt.md) - Transcribe audio
- [Gmail](../gmail/index.md) - Send audio files via email
