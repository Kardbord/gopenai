package main

import (
	"fmt"
	"os"

	"github.com/Kardbord/gopenai/audio"
	"github.com/Kardbord/gopenai/authentication"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

const (
	model             = "whisper-1"
	transcriptionFile = "transcription.m4a"
	translationFile   = "translation.m4a"
)

func main() {
	{ // Transcription
		fmt.Printf("Sending transcription request for file %s...\n", transcriptionFile)
		r, err := audio.MakeTranscriptionRequest(&audio.TranscriptionRequest{
			File:     transcriptionFile,
			Model:    model,
			Language: "en",
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

	{ // Speech
		const s string = "The quick brown fox jumps over the lazy dog."
		fmt.Printf("Sending speech creation request for \"%s\"\n", s)
		resp, err := audio.MakeSpeechRequest(&audio.SpeechRequest{
			Model:          "tts-1",
			Input:          s,
			Voice:          audio.VoiceNova,
			ResponseFormat: audio.SpeechFormatMp3,
		}, nil)
		if err != nil {
			fmt.Printf("Error with speech creation request: %s\n", err)
		}
		if len(resp) == 0 {
			fmt.Println("No TTS audio returned. :(")
		} else {
			err = os.WriteFile(fmt.Sprintf("speech-creation.%s", audio.SpeechFormatMp3), resp, 0644)
			if err != nil {
				fmt.Printf("Error writing %s to disk: %s\n", audio.SpeechFormatMp3, err)
			}
		}
	}
}
