# Speech-to-Text hello world

Transcribes a short audio file with the Cloud Speech-to-Text API.

```bash
go run . -c service-account.json -f dreams-of-the-future.flac
```

The sample says "I like the dreams of the future better than the history of
the past." It was synthesized with the Cloud Text-to-Speech API using the
`en-US-Wavenet-A` voice (see `texttospeech/v1beta1/examples/quickstart`), so
the repository holds no third-party recording. `dreams-of-the-future.mp3` is
the synthesized file; the `.flac` is the same audio converted to 16 kHz mono:

```bash
ffmpeg -i dreams-of-the-future.mp3 -ac 1 -ar 16000 dreams-of-the-future.flac
```
