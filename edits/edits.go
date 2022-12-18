// Package edits provides bindings for the [edits] [endpoint].
// Given a prompt and an instruction, the model will return
// an edited version of the prompt.
//
// [edits]: https://beta.openai.com/docs/api-reference/edits
// [endpoint]: https://api.openai.com/v1/edits
package edits

import (
	"errors"
	"net/http"

	"github.com/TannerKvarfordt/gopenai/common"
)

const Endpoint = common.BaseURL + "edits"

// Request structure for the edits API endpoint.
type Request struct {
	// ID of the model to use. You can use the List models API
	// to see all of your available models, or see our Model
	// overview for descriptions of them.
	Model string `json:"model"`

	// The input text to use as a starting point for the edit.
	Input string `json:"input"`

	// The instruction that tells the model how to edit the prompt.
	Instruction string `json:"instruction"`

	// How many edits to generate for the input and instruction.
	N *int64 `json:"n"`

	// What sampling temperature to use. Higher values means the model
	// will take more risks. Try 0.9 for more creative applications,
	// and 0 (argmax sampling) for ones with a well-defined answer.
	// We generally recommend altering this or top_p but not both.
	Temperature *float64 `json:"temperature"`

	// An alternative to sampling with temperature, called nucleus sampling,
	// where the model considers the results of the tokens with top_p
	// probability mass. So 0.1 means only the tokens comprising the top 10%
	// probability mass are considered.
	// We generally recommend altering this or temperature but not both.
	TopP *float64 `json:"top_p"`
}

// Response structure for the edits API endpoint.
type Response struct {
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Text  string `json:"text"`
		Index int64  `json:"index"`
	} `json:"choices"`
	Usage common.ResponseUsage  `json:"usage"`
	Error *common.ResponseError `json:"error,omitempty"`
}

// Make an edits request.
func MakeRequest(request *Request, organizationID *string) (*Response, error) {
	r, err := common.MakeRequest[Request, Response](request, Endpoint, http.MethodPost, organizationID)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, errors.New("nil response received")
	}
	if r.Error != nil {
		return r, r.Error
	}
	if len(r.Choices) == 0 {
		return r, errors.New("no choices in response")
	}
	return r, nil
}
