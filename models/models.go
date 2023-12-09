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
	"errors"
	"net/http"

	"github.com/Kardbord/gopenai/common"
)

// The models API endpoint.
const Endpoint = common.BaseURL + "models"

// Response structure for a Retrieve Model request.
type ModelResponse struct {
	ID      string `json:"id"`
	Created uint64 `json:"created"`
	Object  string `json:"object"`
	OwnedBy string `json:"owned_by"`

	// Deprecated: No longer listed in the API docs.
	Root string `json:"root"`
	// Deprecated: No longer listed in the API docs.
	Parent *string `json:"parent"`

	Error *common.ResponseError `json:"error,omitempty"`

	// Deprecated: No longer listed in the API docs.
	//
	// The values of each permission object (aka, map)
	// in this list are non-homogeneous. Generally,
	// they are strings, integers, or booleans, but
	// it very much depends on the individual model.
	Permission []map[string]any `json:"permission"`
}

// Response structure for a List Models request.
type ListModelsResponse struct {
	Object string                `json:"object"`
	Data   []ModelResponse       `json:"data"`
	Error  *common.ResponseError `json:"error,omitempty"`
}

// Lists the currently available models, and provides basic information about each one such as the owner and availability.
func MakeListModelsRequest(organizationID *string) (*ListModelsResponse, error) {
	r, err := common.MakeRequest[any, ListModelsResponse](nil, Endpoint, http.MethodGet, organizationID)
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

// Retrieves a model instance, providing basic information about the model such as the owner and permissioning.
func MakeRetrieveModelRequest(model string, organizationID *string) (*ModelResponse, error) {
	r, err := common.MakeRequest[any, ModelResponse](nil, Endpoint+"/"+model, http.MethodGet, organizationID)
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
