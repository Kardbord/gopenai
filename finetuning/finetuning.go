// Package fine-tuning provides bindings for the [finetunes] [endpoint].
// Manage fine-tuning jobs to tailor a model to your specific training data.
// Related guide: [Fine-tune models].
//
// [finetunes]: https://platform.openai.com/docs/api-reference/fine-tuning
// [endpoint]: https://api.openai.com/v1/fine_tuning/jobs
// [Fine-tune models]: https://platform.openai.com/docs/guides/fine-tuning
package finetuning

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Kardbord/gopenai/common"
	"github.com/Kardbord/gopenai/models"
)

const Endpoint = common.BaseURL + "fine_tuning/jobs"

type HyperParameters struct {
	// Number of examples in each batch. A larger batch size means that model
	// parameters are updated less frequently, but with lower variance.
	BatchSize uint64 `json:"batch_size"`

	// Scaling factor for the learning rate. A smaller learning rate may be
	// useful to avoid overfitting.
	LearningRateMultiplier float64 `json:"learning_rate_multiplier"`

	// The number of epochs to train the model for. An epoch refers to one
	// full cycle through the training dataset.
	NEpochs uint64 `json:"n_epochs"`
}

// Request structure for the "create" fine-tune endpoint.
type CreationRequest struct {
	// The name of the model to fine-tune. You can select one of the supported [models].
	//
	// [models]: https://platform.openai.com/docs/models/overview
	Model string `json:"model,omitempty"`

	// The ID of an uploaded file that contains training data.
	//
	// See [upload file] for how to upload a file.
	//
	// Your dataset must be formatted as a JSONL file. Additionally, you must upload your file with the purpose fine-tune.
	//
	// See the [fine-tuning guide] for more details.
	//
	// [upload file]: https://beta.openai.com/docs/api-reference/files/upload
	// [fine-tuning guide]: https://beta.openai.com/docs/guides/fine-tuning/creating-training-data
	TrainingFile string `json:"training_file,omitempty"`

	// The hyperparameters used for the fine-tuning job.
	Hyperparams HyperParameters `json:"hyperparameters"`

	// A string of up to 18 characters that will be added to your fine-tuned model name.
	//
	// For example, a suffix of "custom-model-name" would produce a model name like
	// ft:gpt-3.5-turbo:openai:custom-model-name:7p4lURel.
	Suffix *string `json:"suffix,omitempty"`

	// The ID of an uploaded file that contains validation data.
	//
	// If you provide this file, the data is used to generate validation metrics periodically
	// during fine-tuning. These metrics can be viewed in the fine-tuning results file. The
	// same data should not be present in both train and validation files.
	//
	// Your dataset must be formatted as a JSONL file. You must upload your file with the
	// purpose fine-tune.
	//
	// See the [fine-tuning guide] for more details.
	//
	// [fine-tuning results]: https://beta.openai.com/docs/guides/fine-tuning/analyzing-your-fine-tuned-model
	// [fine-tuning guide]: https://beta.openai.com/docs/guides/fine-tuning/creating-training-data
	ValidationFile *string `json:"validation_file,omitempty"`
}

type FineTuneEvent struct {
	ID        string `json:"id"`
	CreatedAt uint64 `json:"created_at"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Object    string `json:"object"`
	Type      string `json:"type"`
}

// Response structure for the "create" fine-tune endpoint.
type FineTune struct {
	ID             string                `json:"id"`
	CreatedAt      uint64                `json:"created_at"`
	Error          *common.ResponseError `json:"error,omitempty"`
	FineTunedModel string                `json:"fine_tuned_model"`
	FinishedAt     uint64                `json:"finished_at"`
	Hyperparams    HyperParameters       `json:"hyperparameters"`
	Model          string                `json:"model"`
	Object         string                `json:"object"`
	OrganizationID string                `json:"organization_id"`
	ResultFiles    []string              `json:"result_files"`
	Status         string                `json:"status"`
	TrainedTokens  uint64                `json:"trained_tokens"`
	TrainingFile   string                `json:"training_file"`
	ValidationFile string                `json:"validation_file"`
}

// Creates a job that fine-tunes a specified model from a given dataset.
//
// Response includes details of the enqueued job including job status and
// the name of the fine-tuned models once complete.
//
// [Learn more about Fine-tuning]
//
// [Learn more about Fine-tuning]: https://beta.openai.com/docs/guides/fine-tuning
func MakeCreationRequest(request *CreationRequest, organizationID *string) (*FineTune, error) {
	r, err := common.MakeRequest[CreationRequest, FineTune](request, Endpoint, http.MethodPost, organizationID)
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

type ListResponse struct {
	Object  string                `json:"object"`
	Data    []FineTune            `json:"data"`
	Error   *common.ResponseError `json:"error,omitempty"`
	HasMore bool                  `json:"has_more"`
}

// List your organization's fine-tuning jobs
func MakeListRequest(limit *uint64, after, organizationID *string) (*ListResponse, error) {
	endpoint := Endpoint
	if after != nil && limit != nil {
		endpoint = fmt.Sprintf("%s?after=%s&limit=%d", endpoint, *after, *limit)
	} else if after != nil {
		endpoint = fmt.Sprintf("%s?after=%s", endpoint, *after)
	} else if limit != nil {
		endpoint = fmt.Sprintf("%s?limit=%d", endpoint, *limit)
	}
	r, err := common.MakeRequest[any, ListResponse](nil, endpoint, http.MethodGet, organizationID)
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

// Gets info about the fine-tune job.
func MakeRetrieveRequest(fineTuneID string, organizationID *string) (*FineTune, error) {
	r, err := common.MakeRequest[any, FineTune](nil, fmt.Sprintf("%s/%s", Endpoint, fineTuneID), http.MethodGet, organizationID)
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

// Immediately cancel a fine-tune job.
func MakeCancelRequest(fineTuneID string, organizationID *string) (*FineTune, error) {
	r, err := common.MakeRequest[any, FineTune](nil, fmt.Sprintf("%s/%s/cancel", Endpoint, fineTuneID), http.MethodPost, organizationID)
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

type ListEventsResponse struct {
	Object  string                `json:"object"`
	Data    []FineTuneEvent       `json:"data"`
	Error   *common.ResponseError `json:"error,omitempty"`
	HasMore bool                  `json:"has_more"`
}

// Get fine-grained status updates for a fine-tune job.
func MakeListEventsRequest(fineTuneID string, limit *uint64, after, organizationID *string) (*ListEventsResponse, error) {
	// TODO: support streaming: https://beta.openai.com/docs/api-reference/fine-tunes/events#fine-tunes/events-stream

	endpoint := fmt.Sprintf("%s/%s/events", Endpoint, fineTuneID)
	if after != nil && limit != nil {
		endpoint = fmt.Sprintf("%s?after=%s&limit=%d", endpoint, *after, *limit)
	} else if after != nil {
		endpoint = fmt.Sprintf("%s?after=%s", endpoint, *after)
	} else if limit != nil {
		endpoint = fmt.Sprintf("%s?limit=%d", endpoint, *limit)
	}

	r, err := common.MakeRequest[any, ListEventsResponse](nil, endpoint, http.MethodGet, organizationID)
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

type DeleteResponse struct {
	ID      int64                 `json:"id"`
	Object  string                `json:"object"`
	Deleted bool                  `json:"deleted"`
	Error   *common.ResponseError `json:"error,omitempty"`
}

// Delete a fine-tuned model. You must have the Owner role in your organization.
func MakeDeleteRequest(fineTuneModel string, organizationID *string) (*DeleteResponse, error) {
	r, err := common.MakeRequest[any, DeleteResponse](nil, fmt.Sprintf("%s/%s", models.Endpoint, fineTuneModel), http.MethodDelete, organizationID)
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
