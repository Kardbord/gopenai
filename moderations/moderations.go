// Package moderations provides bindings for the [moderations] [endpoint].
//
// [moderations]: https://beta.openai.com/docs/api-reference/moderations
// [endpoint]: https://api.openai.com/v1/moderations
package moderations

import (
	"net/http"

	"github.com/TannerKvarfordt/gopenai/common"
)

// The moderations API endpoint.
const Endpoint = common.BaseURL + "moderations"

const (
	// The name of the stable moderation model.
	ModelStable = "text-moderation-stable"

	// The name of the latest moderation model.
	ModelLatest = "text-moderation-latest"
)

// The request structure for moderation requests.
type Request struct {
	// The input text to classify.
	Input string `json:"input"`

	// Two content moderations models are available: text-moderation-stable and text-moderation-latest.
	// The default is text-moderation-latest which will be automatically upgraded over time.
	// This ensures you are always using our most accurate model. If you use text-moderation-stable,
	// we will provide advanced notice before updating the model. Accuracy of text-moderation-stable
	// may be slightly lower than for text-moderation-latest.
	Model string `json:"model,omitempty"`
}

// The response structure for moderation endpoint responses.
type Response struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Results []struct {
		// Set to true if the model classifies the content as violating OpenAI's content policy, false otherwise.
		Flagged bool `json:"flagged"`

		// Contains a dictionary of per-category binary content policy violation flags. For each category, the value is true if the model flags the corresponding category as violated, false otherwise.
		Categories map[string]bool `json:"categories"`

		// Contains a dictionary of per-category raw scores output by the model, denoting the model's confidence that the input violates the OpenAI's policy for the category. The value is between 0 and 1, where higher values denote higher confidence. The scores should not be interpreted as probabilities.
		CategoryScores map[string]float64 `json:"category_scores"`
	} `json:"results"`

	Error *common.ResponseError `json:"error,omitempty"`
}

func MakeRequest(request *Request, organizationID *string) (*Response, error) {
	response := new(Response)
	err := common.MakeRequest(request, response, Endpoint, http.MethodPost, organizationID)
	return response, err
}
