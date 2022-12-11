// Package moderations provides bindings for the [moderations] [endpoint].
//
// [moderations]: https://beta.openai.com/docs/api-reference/moderations
// [endpoint]: https://api.openai.com/v1/moderations
package moderations

import (
	"net/http"

	openaicommon "github.com/TannerKvarfordt/gopenai/openai-common"
)

// The moderations API endpoint.
const Endpoint = openaicommon.BaseURL + "/moderations"

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
		Categories struct {
			// Content that expresses, incites, or promotes hate based on race, gender, ethnicity, religion, nationality, sexual orientation, disability status, or caste.
			Hate bool `json:"hate"`

			// Hateful content that also includes violence or serious harm towards the targeted group.
			HateThreatening bool `json:"hate/threatening"`

			// Content that promotes, encourages, or depicts acts of self-harm, such as suicide, cutting, and eating disorders.
			SelfHarm bool `json:"self-harm"`

			// Content meant to arouse sexual excitement, such as the description of sexual activity, or that promotes sexual services (excluding sex education and wellness).
			Sexual bool `json:"sexual"`

			// Sexual content that includes an individual who is under 18 years old.
			SexualMinors bool `json:"sexual/minors"`

			// Content that promotes or glorifies violence or celebrates the suffering or humiliation of others.
			Violence bool `json:"violence"`

			// Violent content that depicts death, violence, or serious physical injury in extreme graphic detail.
			ViolenceGraphic bool `json:"violence/graphic"`
		} `json:"categories"`

		// Contains a dictionary of per-category raw scores output by the model, denoting the model's confidence that the input violates the OpenAI's policy for the category. The value is between 0 and 1, where higher values denote higher confidence. The scores should not be interpreted as probabilities.
		CategoryScores struct {
			Hate            float64 `json:"hate"`
			HateThreatening float64 `json:"hate/threatening"`
			SelfHarm        float64 `json:"self-harm"`
			Sexual          float64 `json:"sexual"`
			SexualMinors    float64 `json:"sexual/minors"`
			Violence        float64 `json:"violence"`
			ViolenceGraphic float64 `json:"violence/graphic"`
		} `json:"category_scores"`
	} `json:"results"`
}

func MakeModerationRequest(request *Request, organizationID *string) (*Response, error) {
	response := new(Response)
	err := openaicommon.MakeRequest(request, response, Endpoint, http.MethodPost, organizationID)
	if err != nil {
		return nil, err
	}
	return response, nil
}
