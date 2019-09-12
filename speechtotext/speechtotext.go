package speechtotext

import (
	"fmt"
	"io/ioutil"
	"strings"

	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

/*

This is based on the following Quickstart

https://github.com/GoogleCloudPlatform/golang-samples/blob/master/dlp/dlp_quickstart/quickstart.go

*/

func NewRecognitionAudio(data []byte) *speechpb.RecognitionAudio {
	return &speechpb.RecognitionAudio{
		AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
	}
}

func NewRecognitionAudioFile(file string) (*speechpb.RecognitionAudio, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return NewRecognitionAudio(data), nil
}

func RecognizeResponseTextFirst(resp *speechpb.RecognizeResponse, threshold float32) (string, error) {
	//highestThreshold := 0.0
	//highestThresholdTranscript := ""
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			trimmed := strings.TrimSpace(alt.Transcript)
			if len(trimmed) > 0 {
				if threshold == 0 ||
					(threshold > 0 && alt.Confidence >= threshold) {
					return trimmed, nil
				}
			}
		}
	}
	return "", fmt.Errorf("No responses found")
}
