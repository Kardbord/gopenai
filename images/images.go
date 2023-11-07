// Package images provides bindings for the [images] [endpoint].
// Given a prompt and/or an input image, the model will generate a new image.
// Related guide: [Image Generation].
//
// [images]: https://beta.openai.com/docs/api-reference/images
// [endpoint]: https://api.openai.com/v1/images
// [Image Generation]: https://beta.openai.com/docs/guides/images/image-generation-beta
package images

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/TannerKvarfordt/gopenai/common"
	"github.com/TannerKvarfordt/gopenai/moderations"
)

const (
	BaseEndpoint      = common.BaseURL + "images/"
	CreateEndpoint    = BaseEndpoint + "generations"
	EditEndpoint      = BaseEndpoint + "edits"
	VariationEndpoint = BaseEndpoint + "variations"
)

const (
	Dalle2SmallImage  = "256x256"
	Dalle2MediumImage = "512x512"
	Dalle2LargeImage  = "1024x1024"

	Dalle3SquareImage    = "1024x1024"
	Dalle3LandscapeImage = "1792x1024"
	Dalle3PortraitImage  = "1024x1792"

	// Deprecated: Use Dalle2SmallImage instead.
	SmallImage = Dalle2SmallImage
	// Deprecated: Use Dalle2MediumImage instead.
	MediumImage = Dalle2MediumImage
	// Deprecated: Use Dalle2LargeImage instead.
	LargeImage = Dalle2LargeImage
)

const (
	ResponseFormatURL     = "url"
	ResponseFormatB64JSON = "b64_json"
)

const (
	ModelDalle2 = "dall-e-2"
	ModelDalle3 = "dall-e-3"
)

const (
	QualityStandard = "standard"
	QualityHD       = "hd"
)

const (
	StyleVivid   = "vivid"
	StyleNatural = "natural"
)

// Response structure for the image API endpoint.
type Response struct {
	Created uint64 `json:"created"`
	Data    []struct {
		URL     string `json:"url"`
		B64JSON string `json:"b64_json"`
	}
	Error *common.ResponseError `json:"error,omitempty"`
}

// Request structure for the image creation API endpoint.
type CreationRequest struct {
	// A text description of the desired image(s).
	// The maximum length is 1000 characters for dall-e-2 and 4000 characters for dall-e-3.
	Prompt string `json:"prompt,omitempty"`

	// The model to use for image generation.
	Model string `json:"model,omitempty"`

	// The number of images to generate. Must be between 1 and 10.
	// For dall-e-3, only n=1 is supported.
	N *uint64 `json:"n,omitempty"`

	// The quality of the image that will be generated.
	// "hd" creates images with finer details and greater consistency across the image.
	// This param is only supported for dall-e-3.
	Quality string `json:"quality,omitempty"`

	// The format in which the generated images are returned. Must be one of url or b64_json.
	ResponseFormat string `json:"response_format,omitempty"`

	// The size of the generated images.
	// Must be one of 256x256, 512x512, or 1024x1024 for dall-e-2.
	// Must be one of 1024x1024, 1792x1024, or 1024x1792 for dall-e-3 models.
	Size string `json:"size,omitempty"`

	// The style of the generated images. Must be one of vivid or natural.
	// Vivid causes the model to lean towards generating hyper-real and dramatic images.
	// Natural causes the model to produce more natural, less hyper-real looking images.
	// This param is only supported for dall-e-3.
	Style string `json:"style,omitempty"`

	// A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse.
	User string `json:"user,omitempty"`
}

// Creates an image given a prompt.
func MakeCreationRequest(request *CreationRequest, organizationID *string) (*Response, error) {
	r, err := common.MakeRequest[CreationRequest, Response](request, CreateEndpoint, http.MethodPost, organizationID)
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
func MakeModeratedCreationRequest(request *CreationRequest, organizationID *string) (*Response, *moderations.Response, error) {
	modr, err := moderations.MakeModeratedRequest(&moderations.Request{
		Input: []string{request.Prompt},
		Model: moderations.ModelLatest,
	}, organizationID)
	if err != nil {
		return nil, modr, err
	}

	r, err := MakeCreationRequest(request, organizationID)
	if err != nil {
		return nil, modr, err
	}
	return r, modr, nil
}

// Request structure for the image editing API endpoint.
type EditRequest struct {
	// The image to edit. Must be a valid PNG file, less than 4MB, and square.
	// If mask is not provided, image must have transparency, which will be
	// used as the mask.
	Image string `json:"image,omitempty"`

	// The name of the image, including its extension, but not including
	// any path information.
	ImageName string `json:"-"`

	// An additional image whose fully transparent areas (e.g. where alpha is zero)
	// indicate where image should be edited. Must be a valid PNG file, less than 4MB,
	// and have the same dimensions as image.
	Mask string `json:"mask,omitempty"`

	// The name of the mask, including its extension, but not including any
	// path information.
	MaskName string `json:"-"`

	// A text description of the desired image(s). The maximum length is 1000 characters.
	Prompt string `json:"prompt,omitempty"`

	// The number of images to generate. Must be between 1 and 10.
	N *uint64 `json:"n,omitempty"`

	// The size of the generated images. Must be one of 256x256, 512x512, or 1024x1024.
	Size string `json:"size,omitempty"`

	// The format in which the generated images are returned. Must be one of url or b64_json.
	ResponseFormat string `json:"response_format,omitempty"`

	// A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse.
	User string `json:"user,omitempty"`
}

// Creates an edited or extended image given an original image and a prompt.
func MakeEditRequest(request *EditRequest, organizationID *string) (*Response, error) {
	if request == nil {
		return nil, errors.New("nil request provided")
	}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	if len(request.Prompt) > 0 {
		err := common.CreateFormField("prompt", request.Prompt, writer)
		if err != nil {
			return nil, err
		}
	}

	var err error
	if request.N != nil {
		err = common.CreateFormField("n", request.N, writer)
		if err != nil {
			return nil, err
		}
	}

	if len(request.Size) > 0 {
		err = common.CreateFormField("size", request.Size, writer)
		if err != nil {
			return nil, err
		}
	}

	if len(request.ResponseFormat) > 0 {
		err = common.CreateFormField("response_format", request.ResponseFormat, writer)
		if err != nil {
			return nil, err
		}
	}

	if len(request.User) > 0 {
		err = common.CreateFormField("user", request.User, writer)
		if err != nil {
			return nil, err
		}
	}

	if len(request.Image) > 0 {
		err = common.CreateFormFile("image", request.ImageName, request.Image, writer)
		if err != nil {
			return nil, err
		}
	}

	if len(request.Mask) > 0 {
		err = common.CreateFormFile("mask", request.MaskName, request.Mask, writer)
		if err != nil {
			return nil, err
		}
	}

	writer.Close()
	r, err := common.MakeRequestWithForm[Response](buf, EditEndpoint, http.MethodPost, writer.FormDataContentType(), organizationID)
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
func MakeModeratedRequest(request *EditRequest, organizationID *string) (*Response, *moderations.Response, error) {
	modr, err := moderations.MakeModeratedRequest(&moderations.Request{
		Input: []string{request.Prompt},
		Model: moderations.ModelLatest,
	}, organizationID)
	if err != nil {
		return nil, modr, err
	}

	r, err := MakeEditRequest(request, organizationID)
	if err != nil {
		return nil, modr, err
	}
	return r, modr, nil
}

// Request structure for the image variations API endpoint.
type VariationRequest struct {
	// The image to use as the basis for the variation(s). Must be a valid PNG file, less than 4MB, and square.
	Image string `json:"image,omitempty"`

	// The name of the image, including its extension, but not including
	// any path information.
	ImageName string `json:"-"`

	// The number of images to generate. Must be between 1 and 10.
	N *uint64 `json:"n,omitempty"`

	// The size of the generated images. Must be one of 256x256, 512x512, or 1024x1024.
	Size string `json:"size,omitempty"`

	// The format in which the generated images are returned. Must be one of url or b64_json.
	ResponseFormat string `json:"response_format,omitempty"`

	// A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse.
	User string `json:"user,omitempty"`
}

// Creates a variation of a given image.
func MakeVariationRequest(request *VariationRequest, organizationID *string) (*Response, error) {
	if request == nil {
		return nil, errors.New("nil request provided")
	}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	var err error
	if request.N != nil {
		err = common.CreateFormField("n", request.N, writer)
		if err != nil {
			return nil, err
		}
	}

	if len(request.Size) > 0 {
		err = common.CreateFormField("size", request.Size, writer)
		if err != nil {
			return nil, err
		}
	}

	if len(request.ResponseFormat) > 0 {
		err = common.CreateFormField("response_format", request.ResponseFormat, writer)
		if err != nil {
			return nil, err
		}
	}

	if len(request.User) > 0 {
		err = common.CreateFormField("user", request.User, writer)
		if err != nil {
			return nil, err
		}
	}

	if len(request.Image) > 0 {
		err = common.CreateFormFile("image", request.ImageName, request.Image, writer)
		if err != nil {
			return nil, err
		}
	}

	writer.Close()
	r, err := common.MakeRequestWithForm[Response](buf, VariationEndpoint, http.MethodPost, writer.FormDataContentType(), organizationID)
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
