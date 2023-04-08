package audio_test

import (
	"os"
	"testing"

	"github.com/TannerKvarfordt/gopenai/audio"
	"github.com/TannerKvarfordt/gopenai/authentication"
)

const (
	OpenAITokenEnv        = "OPENAI_API_KEY"
	transcriptionFilePath = "./test_files/transcription.m4a"
	translationFilePath   = "./test_files/translation.m4a"
	model                 = "whisper-1"
)

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func TestTranscription(t *testing.T) {
	resp, err := audio.MakeTranscriptionRequest(&audio.TranscriptionRequest{
		File:  transcriptionFilePath,
		Model: model,
	}, nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	if resp.Text == "" {
		t.Fatal("no transcription returned")
		return
	}
}

func TestTranslation(t *testing.T) {
	resp, err := audio.MakeTranslationRequest(&audio.TranslationRequest{
		File:  translationFilePath,
		Model: model,
	}, nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	if resp.Text == "" {
		t.Fatal("no translation returned")
		return
	}
}
