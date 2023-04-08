// Package audio provides bindings for the [audio] [endpoint].
// Converts audio into text.
//
// [endpoint]: https://api.openai.com/v1/audio/transcriptions
//
// [chat]: https://platform.openai.com/docs/api-reference/audio
package audio

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/TannerKvarfordt/gopenai/common"
)

const (
	BaseEndpoint         = common.BaseURL + "audio/"
	TransciptionEndpoint = BaseEndpoint + "transcriptions"
	TranslationEndpoint  = BaseEndpoint + "translations"
)

type ResponseFormat = string

const (
	// TODO: Support non-json return formats.
	JSONResponseFormat = "json"
	//TextResponseFormat        = "text"
	//SRTResponseFormat         = "srt"
	//VerboseJSONResponseFormat = "verbose_json"
	//VTTResponseFormat         = "vtt"
)

// Request structure for the transcription endpoint.
type TranscriptionRequest struct {
	// The audio file to transcribe, in one of these formats:
	// mp3, mp4, mpeg, mpga, m4a, wav, or webm.
	// This can be a file path or a URL.
	File string `json:"file"`

	// ID of the model to use. You can use the List models API
	// to see all of your available models, or see our Model
	// overview for descriptions of them.
	Model string `json:"model"`

	// An optional text to guide the model's style or continue a
	// previous audio segment. The prompt should match the audio language.
	Prompt string `json:"prompt,omitempty"`

	// The format of the transcript output, in one of these options:
	// json, text, srt, verbose_json, or vtt.
	ResponseFormat ResponseFormat `json:"response_format,omitempty"`

	// The sampling temperature, between 0 and 1. Higher values like 0.8 will
	// make the output more random, while lower values like 0.2 will make it
	// more focused and deterministic. If set to 0, the model will use log
	// probability to automatically increase the temperature until certain
	// thresholds are hit.
	Temperature *float64 `json:"temperature,omitempty"`

	// The language of the input audio. Supplying the input language in
	// ISO-639-1 format will improve accuracy and latency.
	Language string `json:"language,omitempty"`
}

// Request structure for the Translations endpoint.
type TranslationRequest struct {
	// The audio file to transcribe, in one of these formats:
	// mp3, mp4, mpeg, mpga, m4a, wav, or webm.
	// This can be a file path or a URL.
	File string `json:"file"`

	// ID of the model to use. You can use the List models API
	// to see all of your available models, or see our Model
	// overview for descriptions of them.
	Model string `json:"model"`

	// An optional text to guide the model's style or continue a
	// previous audio segment. The prompt should be in English.
	Prompt string `json:"prompt,omitempty"`

	// The format of the transcript output, in one of these options:
	// json, text, srt, verbose_json, or vtt.
	ResponseFormat ResponseFormat `json:"response_format,omitempty"`

	// The sampling temperature, between 0 and 1. Higher values like 0.8 will
	// make the output more random, while lower values like 0.2 will make it
	// more focused and deterministic. If set to 0, the model will use log
	// probability to automatically increase the temperature until certain
	// thresholds are hit.
	Temperature *float64 `json:"temperature,omitempty"`
}

// Response structure for both Transcription and
// Translation requests.
type Response struct {
	Text  string                `json:"text"`
	Usage common.ResponseUsage  `json:"usage"`
	Error *common.ResponseError `json:"error,omitempty"`
}

func MakeTranscriptionRequest(request *TranscriptionRequest, organizationID *string) (*Response, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	err := common.CreateFormField("model", request.Model, writer)
	if err != nil {
		return nil, err
	}

	err = common.CreateFormFile("file", filepath.Base(request.File), request.File, writer)
	if err != nil {
		return nil, err
	}
	writer.Close()
	r, err := common.MakeRequestWithForm[Response](buf, TransciptionEndpoint, http.MethodPost, writer.FormDataContentType(), organizationID)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, errors.New("nil response received")
	}
	if r.Error != nil {
		return r, r.Error
	}
	return r, nil
}

func MakeTranslationRequest(request *TranslationRequest, organizationID *string) (*Response, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	err := common.CreateFormField("model", request.Model, writer)
	if err != nil {
		return nil, err
	}

	err = common.CreateFormFile("file", filepath.Base(request.File), request.File, writer)
	if err != nil {
		return nil, err
	}
	writer.Close()
	r, err := common.MakeRequestWithForm[Response](buf, TranslationEndpoint, http.MethodPost, writer.FormDataContentType(), organizationID)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, errors.New("nil response received")
	}
	if r.Error != nil {
		return r, r.Error
	}
	return r, nil
}
