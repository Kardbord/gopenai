package main

import (
	"fmt"
	"os"

	"github.com/TannerKvarfordt/gopenai/audio"
	"github.com/TannerKvarfordt/gopenai/authentication"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

const (
	model             = "whisper-1"
	transcriptionFile = "./transcription.m4a"
	translationFile   = "./translation.m4a"
)

func main() {
	{ // Transcription
		fmt.Printf("Sending transcription request for file %s...\n", transcriptionFile)
		r, err := audio.MakeTranscriptionRequest(&audio.TranscriptionRequest{
			File:  transcriptionFile,
			Model: model,
		}, nil)
		if err != nil {
			fmt.Printf("Error with transcription request: %s\n", err)
		} else {
			fmt.Printf("Transcribed audio: %s\n", r.Text)
		}
	}

	{ // Translation
		fmt.Printf("Sending translation request for file %s...\n", translationFile)
		r, err := audio.MakeTranslationRequest(&audio.TranslationRequest{
			File:  translationFile,
			Model: model,
		}, nil)
		if err != nil {
			fmt.Printf("Error with translation request: %s\n", err)
		} else {
			fmt.Printf("Translated audio: %s\n", r.Text)
		}
	}
}
