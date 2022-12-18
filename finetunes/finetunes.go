// Package finetunes provides bindings for the [finetunes] [endpoint].
// Manage fine-tuning jobs to tailor a model to your specific training data.
// Related guide: [Fine-tune models].
//
// [finetunes]: https://beta.openai.com/docs/api-reference/finetunes
// [endpoint]: https://api.openai.com/v1/finetunes
// [Fine-tune models]: https://beta.openai.com/docs/guides/fine-tuning
package finetunes

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/TannerKvarfordt/gopenai/common"
	"github.com/TannerKvarfordt/gopenai/files"
	"github.com/TannerKvarfordt/gopenai/models"
)

const Endpoint = common.BaseURL + "fine-tunes"

// Request structure for the "create" fine-tune endpoint.
type CreationRequest struct {

	// The ID of an uploaded file that contains training data.
	//
	// See [upload file] for how to upload a file.
	//
	// Your dataset must be formatted as a JSONL file, where each training example
	// is a JSON object with the keys "prompt" and "completion". Additionally, you
	// must upload your file with the purpose fine-tune.
	//
	// See the [fine-tuning guide] for more details.
	//
	// [upload file]: https://beta.openai.com/docs/api-reference/files/upload
	// [fine-tuning guide]: https://beta.openai.com/docs/guides/fine-tuning/creating-training-data
	TrainingFile string `json:"training_file,omitempty"`
	// The ID of an uploaded file that contains validation data.
	//
	// If you provide this file, the data is used to generate validation metrics periodically
	// during fine-tuning. These metrics can be viewed in the [fine-tuning results] file.
	// Your train and validation data should be mutually exclusive.
	//
	// Your dataset must be formatted as a JSONL file, where each validation example is a
	// JSON object with the keys "prompt" and "completion". Additionally, you must upload
	// your file with the purpose fine-tune.
	//
	// See the [fine-tuning guide] for more details.
	//
	// [fine-tuning results]: https://beta.openai.com/docs/guides/fine-tuning/analyzing-your-fine-tuned-model
	// [fine-tuning guide]: https://beta.openai.com/docs/guides/fine-tuning/creating-training-data
	ValidationFile *string `json:"validation_file,omitempty"`
	// The name of the base model to fine-tune. You can select one of "ada", "babbage", "curie",
	// "davinci", or a fine-tuned model created after 2022-04-21. To learn more about these models,
	// see the [Models] documentation.
	//
	// Defaults to "curie".
	//
	// [Models]: https://beta.openai.com/docs/models
	Model *string `json:"model,omitempty"`
	// The number of epochs to train the model for. An epoch refers to one full cycle through the
	// training dataset.
	NEpochs *uint64 `json:"n_epochs,omitempty"`
	// The batch size to use for training. The batch size is the number of training examples used
	// to train a single forward and backward pass.
	//
	// By default, the batch size will be dynamically configured to be ~0.2% of the number of
	// examples in the training set, capped at 256 - in general, we've found that larger batch
	// sizes tend to work better for larger datasets.
	BatchSize *uint64 `json:"batch_size,omitempty"`
	// The learning rate multiplier to use for training. The fine-tuning learning rate is the
	// original learning rate used for pretraining multiplied by this value.
	//
	// By default, the learning rate multiplier is the 0.05, 0.1, or 0.2 depending on final
	// batch_size (larger learning rates tend to perform better with larger batch sizes).
	// We recommend experimenting with values in the range 0.02 to 0.2 to see what produces
	// the best results.
	LearningRateMultiplier *float64 `json:"learning_rate_multiplier,omitempty"`
	// The weight to use for loss on the prompt tokens. This controls how much the model tries
	// to learn to generate the prompt (as compared to the completion which always has a weight
	// of 1.0), and can add a stabilizing effect to training when completions are short.
	//
	// If prompts are extremely long (relative to completions), it may make sense to reduce
	// this weight so as to avoid over-prioritizing learning the prompt.
	PromptLossWeight *float64 `json:"prompt_loss_weight,omitempty"`
	// If set, we calculate classification-specific metrics such as accuracy and F-1 score using
	// the validation set at the end of every epoch. These metrics can be viewed in the results file.
	//
	// In order to compute classification metrics, you must provide a validation_file. Additionally,
	// you must specify classification_n_classes for multiclass classification or classification_positive_class
	// for binary classification.
	ComputeClassificationMetrics bool `json:"compute_classification_metrics,omitempty"`
	// The number of classes in a classification task. This parameter is required for multiclass classification.
	ClassificationNClasses *uint64 `json:"classification_n_classes,omitempty"`
	// The positive class in binary classification.
	//
	// This parameter is needed to generate precision, recall, and F1 metrics when doing binary classification.
	ClassificationPositiveClass *string `json:"classification_positive_class,omitempty"`
	// A string of up to 40 characters that will be added to your fine-tuned model name.
	//
	// For example, a suffix of "custom-model-name" would produce a model name like
	// ada:ft-your-org:custom-model-name-2022-02-15-04-21-04.
	Suffix *string `json:"suffix,omitempty"`
	// TODO: Add support for classification_betas: https://beta.openai.com/docs/api-reference/fine-tunes/create#fine-tunes/create-classification_betas
	//ClassificationBetas
}

type FineTuneEvent struct {
	Object    string `json:"object"`
	CreatedAt uint64 `json:"created_at"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// Response structure for the "create" fine-tune endpoint.
type FineTune struct {
	ID             string          `json:"id"`
	Object         string          `json:"object"`
	Model          string          `json:"model"`
	CreatedAt      uint64          `json:"created_at"`
	Events         []FineTuneEvent `json:"events"`
	FineTunedModel string          `json:"fine_tuned_model"`
	Hyperparams    struct {
		BatchSize              uint64  `json:"batch_size"`
		LearningRateMultiplier float64 `json:"learning_rate_multiplier"`
		NEpochs                uint64  `json:"n_epochs"`
		PromptLossWeight       float64 `json:"prompt_loss_weight"`
	} `json:"hyperparams"`
	OrganizationID  string                `json:"organization_id"`
	ResultsFiles    []files.UploadedFile  `json:"results_files"`
	Status          string                `json:"status"`
	ValidationFiles []files.UploadedFile  `json:"validation_files"`
	TrainingFiles   []files.UploadedFile  `json:"training_files"`
	UpdatedAt       uint64                `json:"updated_at"`
	Error           *common.ResponseError `json:"error,omitempty"`
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
	Object string                `json:"object"`
	Data   []FineTune            `json:"data"`
	Error  *common.ResponseError `json:"error,omitempty"`
}

// List your organization's fine-tuning jobs
func MakeListRequest(organizationID *string) (*ListResponse, error) {
	r, err := common.MakeRequest[any, ListResponse](nil, Endpoint, http.MethodGet, organizationID)
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
	Object string                `json:"object"`
	Data   []FineTuneEvent       `json:"data"`
	Error  *common.ResponseError `json:"error,omitempty"`
}

// Get fine-grained status updates for a fine-tune job.
func MakeListEventsRequest(fineTuneID string, organizationID *string) (*ListEventsResponse, error) {
	// TODO: support streaming: https://beta.openai.com/docs/api-reference/fine-tunes/events#fine-tunes/events-stream
	r, err := common.MakeRequest[any, ListEventsResponse](nil, fmt.Sprintf("%s/%s/events", Endpoint, fineTuneID), http.MethodGet, organizationID)
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
