// Package embeddings provides bindings for the [embeddings] [endpoint].
// Get a vector representation of a given input that can be easily consumed
// by machine learning models and algorithms.
//
// [embeddings]: https://beta.openai.com/docs/api-reference/embeddings
// [endpoint]: https://api.openai.com/v1/embeddings
package embeddings

import (
	"errors"
	"net/http"

	"github.com/Kardbord/gopenai/common"
	"github.com/Kardbord/gopenai/moderations"
)

const Endpoint = common.BaseURL + "embeddings"

// Request structure for the embeddings API endpoint.
type Request struct {
	// ID of the model to use. You can use the List models API to see all of
	// your available models, or see our Model overview for descriptions of them.
	Model string `json:"model"`

	// Input text to get embeddings for, encoded as a string or array of tokens.
	// To get embeddings for multiple inputs in a single request, pass an array
	// of strings or array of token arrays. Each input must not exceed 8192 tokens
	// in length.
	Input []string `json:"input"`

	// A unique identifier representing your end-user, which can help OpenAI to
	// monitor and detect abuse.
	User string `json:"user"`

	// The format to return the embeddings in. Can be either float or base64.
	EncodingFormat string `json:"encoding_format,omitempty"`
}

// Response structure for the embeddings API endpoint.
type Response struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
		Index     uint64    `json:"index"`
	} `json:"data"`
	Model string                `json:"model"`
	Usage common.ResponseUsage  `json:"usage"`
	Error *common.ResponseError `json:"error,omitempty"`
}

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
	if len(r.Data) == 0 {
		return r, errors.New("no data in response")
	}
	return r, nil
}

// Runs request inputs through the moderations endpoint prior to making the request.
// Returns a moderations.ModerationFlagError prior to making the request if the
// inputs are flagged by the moderations endpoint.
func MakeModeratedRequest(request *Request, organizationID *string) (*Response, *moderations.Response, error) {
	modr, err := moderations.MakeModeratedRequest(&moderations.Request{
		Input: request.Input,
		Model: moderations.ModelLatest,
	}, organizationID)
	if err != nil {
		return nil, modr, err
	}

	r, err := MakeRequest(request, organizationID)
	if err != nil {
		return nil, modr, err
	}
	return r, modr, nil
}
