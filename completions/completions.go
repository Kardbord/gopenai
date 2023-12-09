// Package completions provides bindings for the [completions] [endpoint].
// Given a prompt, the model will return one or more predicted completions,
// and can also return the probabilities of alternative tokens at each position.
//
// [completions]: https://beta.openai.com/docs/api-reference/completions
// [endpoint]: https://api.openai.com/v1/completions
package completions

import (
	"errors"
	"net/http"

	"github.com/Kardbord/gopenai/common"
	"github.com/Kardbord/gopenai/moderations"
)

// The completions API endpoint.
const Endpoint = common.BaseURL + "completions"

// Request structure for the completions API endpoint.
type Request struct {
	// ID of the model to use. You can use the List models API to see
	// all of your available models, or see our Model overview for
	// descriptions of them.
	Model string `json:"model"`

	// The prompt(s) to generate completions for, encoded as a string,
	// array of strings, array of tokens, or array of token arrays.
	// Note that <|endoftext|> is the document separator that the model
	// sees during training, so if a prompt is not specified the model
	// will generate as if from the beginning of a new document.
	Prompt []string `json:"prompt,omitempty"`

	// The suffix that comes after a completion of inserted text.
	Suffix string `json:"suffix,omitempty"`

	// The maximum number of tokens to generate in the completion.
	// The token count of your prompt plus max_tokens cannot exceed the
	// model's context length. Most models have a context length of 2048
	// tokens (except for the newest models, which support 4096).
	MaxTokens uint64 `json:"max_tokens,omitempty"`

	// What sampling temperature to use. Higher values means the model
	// will take more risks. Try 0.9 for more creative applications,
	// and 0 (argmax sampling) for ones with a well-defined answer.
	// We generally recommend altering this or top_p but not both.
	Temperature *float64 `json:"temperature,omitempty"`

	// An alternative to sampling with temperature, called nucleus sampling,
	// where the model considers the results of the tokens with top_p
	// probability mass. So 0.1 means only the tokens comprising the top 10%
	// probability mass are considered.
	// We generally recommend altering this or temperature but not both.
	TopP *float64 `json:"top_p,omitempty"`

	// How many completions to generate for each prompt.
	// Note: Because this parameter generates many completions, it can quickly
	// consume your token quota. Use carefully and ensure that you have reasonable
	// settings for max_tokens and stop.
	N *uint64 `json:"n,omitempty"`

	// Whether to stream back partial progress. If set, tokens will be sent as
	// data-only server-sent events as they become available, with the stream
	// terminated by a data: [DONE] message.
	// Stream bool `json:"stream,omitempty"` TODO: Add streaming support

	// Include the log probabilities on the logprobs most likely tokens, as well the
	// chosen tokens. For example, if logprobs is 5, the API will return a list of
	// the 5 most likely tokens. The API will always return the logprob of the sampled
	// token, so there may be up to logprobs+1 elements in the response.
	// The maximum value for logprobs is 5. If you need more than this, please contact
	// us through our Help center and describe your use case.
	LogProbs uint64 `json:"logprobs,omitempty"`

	// Echo back the prompt in addition to the completion
	Echo bool `json:"echo,omitempty"`

	// Up to 4 sequences where the API will stop generating further tokens. The returned
	// text will not contain the stop sequence.
	Stop []string `json:"stop,omitempty"`

	// Number between -2.0 and 2.0. Positive values penalize new tokens based on whether
	// they appear in the text so far, increasing the model's likelihood to talk about new topics.
	PresencePenalty float64 `json:"presence_penalty,omitempty"`

	// Number between -2.0 and 2.0. Positive values penalize new tokens based on their existing
	// frequency in the text so far, decreasing the model's likelihood to repeat the same line verbatim.
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`

	// Generates best_of completions server-side and returns the "best" (the one with the highest
	// log probability per token). Results cannot be streamed.
	// When used with n, best_of controls the number of candidate completions and n specifies how
	// many to return – best_of must be greater than n.
	// Note: Because this parameter generates many completions, it can quickly consume your token
	// quota. Use carefully and ensure that you have reasonable settings for max_tokens and stop.
	BestOf *uint64 `json:"best_of,omitempty"`

	// Modify the likelihood of specified tokens appearing in the completion.
	// Accepts a json object that maps tokens (specified by their token ID in the GPT tokenizer)
	// to an associated bias value from -100 to 100. You can use this tokenizer tool (which works
	// for both GPT-2 and GPT-3) to convert text to token IDs. Mathematically, the bias is added to
	// the logits generated by the model prior to sampling. The exact effect will vary per model, but
	// values between -1 and 1 should decrease or increase likelihood of selection; values like
	// -100 or 100 should result in a ban or exclusive selection of the relevant token.
	// As an example, you can pass {"50256": -100} to prevent the <|endoftext|> token from being generated.
	LogitBias map[string]int64 `json:"logit_bias,omitempty"`

	// A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse.
	User string `json:"user,omitempty"`

	// If specified, our system will make a best effort to sample deterministically, such that repeated
	// requests with the same seed and parameters should return the same result. Determinism is not
	// guaranteed, and you should refer to the system_fingerprint response parameter to monitor changes
	// in the backend.
	Seed *int64 `json:"seed,omitempty"`
}

// Response structure for the  completions API endpoint.
type Response struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created uint64 `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string `json:"text"`
		Index        uint64 `json:"index"`
		FinishReason string `json:"finish_reason"`
		LogProbs     struct {
			Tokens        []string             `json:"tokens"`
			TokenLogProbs []float64            `json:"token_logprobs"`
			TopLogProbs   []map[string]float64 `json:"top_logprobs"`
			TextOffset    []uint64             `json:"text_offset"`
		} `json:"logprobs"`
	} `json:"choices"`
	SystemFingerprint string                `json:"system_fingerprint"`
	Usage             common.ResponseUsage  `json:"usage"`
	Error             *common.ResponseError `json:"error,omitempty"`
}

// Make a completions request.
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

// Runs request inputs through the moderations endpoint prior to making the request.
// Returns a moderations.ModerationFlagError prior to making the request if the
// inputs are flagged by the moderations endpoint.
func MakeModeratedRequest(request *Request, organizationID *string) (*Response, *moderations.Response, error) {
	modr, err := moderations.MakeModeratedRequest(&moderations.Request{
		Input: request.Prompt,
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
