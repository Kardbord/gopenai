// Package models provides bindings for the [models] [endpoint].
// List and describe the various models available in the API.
// You can refer to the [Models] documentation to understand what
// models are available and the differences between them.
//
// [Models]: https://beta.openai.com/docs/models
// [models]: https://beta.openai.com/docs/api-reference/models
// [endpoint]: https://api.openai.com/v1/models
package models

import (
	"net/http"

	"github.com/TannerKvarfordt/gopenai/common"
)

// The models API endpoint.
const Endpoint = common.BaseURL + "models"

// Response structure for a Retrieve Model request.
type ModelResponse struct {
	ID      string                `json:"id"`
	Created uint64                `json:"created"`
	OwnedBy string                `json:"owned_by"`
	Root    string                `json:"root"`
	Parent  *string               `json:"parent"`
	Error   *common.ResponseError `json:"error,omitempty"`

	// The values of each permission object (aka, map)
	// in this list are non-homogeneous. Generally,
	// they are strings, integers, or booleans, but
	// it very much depends on the individual model.
	Permission []map[string]any `json:"permission"`
}

// Response structure for a List Models request.
type ListModelsResponse struct {
	Data  []ModelResponse       `json:"data"`
	Error *common.ResponseError `json:"error,omitempty"`
}

// Lists the currently available models, and provides basic information about each one such as the owner and availability.
func MakeListModelsRequest(organizationID *string) (*ListModelsResponse, error) {
	return common.MakeRequest[any, ListModelsResponse](nil, Endpoint, http.MethodGet, organizationID)
}

// Retrieves a model instance, providing basic information about the model such as the owner and permissioning.
func MakeRetrieveModelRequest(model string, organizationID *string) (*ModelResponse, error) {
	return common.MakeRequest[any, ModelResponse](nil, Endpoint+"/"+model, http.MethodGet, organizationID)
}
