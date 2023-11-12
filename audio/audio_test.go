package audio_test

import (
	"errors"
	"os"
	"testing"

	"github.com/TannerKvarfordt/gopenai/audio"
	"github.com/TannerKvarfordt/gopenai/authentication"
	"github.com/TannerKvarfordt/gopenai/common"
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
		File:     transcriptionFilePath,
		Model:    model,
		Language: "en",
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

func TestSpeech(t *testing.T) {
	resp, err := audio.MakeSpeechRequest(&audio.SpeechRequest{
		Model:          "tts-1",
		Input:          "The quick brown fox jumps over the lazy dog.",
		Voice:          audio.VoiceAlloy,
		ResponseFormat: audio.SpeechFormatMp3,
	}, nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	if len(resp) == 0 {
		t.Fatal("No audio returned")
		return
	}
}

func TestInvalidSpeechRequest(t *testing.T) {
	_, err := audio.MakeSpeechRequest(&audio.SpeechRequest{
		Model:          "",
		Input:          "The quick brown fox jumps over the lazy dog.",
		ResponseFormat: audio.SpeechFormatMp3,
	}, nil)
	if err == nil {
		t.Fatal("Expected to receive an invalid request error")
		return
	}
	respErr := new(common.ResponseError)
	if !errors.As(err, &respErr) {
		t.Fatal("Expected error to be of type common.ResponseError")
	}
}
