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
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/TannerKvarfordt/gopenai/common"
)

const (
	BaseEndpoint      = common.BaseURL + "images/"
	CreateEndpoint    = BaseEndpoint + "generations"
	EditEndpoint      = BaseEndpoint + "edits"
	VariationEndpoint = BaseEndpoint + "variations"
)

const (
	SmallImage  string = "256x256"
	MediumImage string = "512x512"
	LargeImage  string = "1024x1024"
)

const (
	ResponseFormatURL     = "url"
	ResponseFormatB64JSON = "b64_json"
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

// Creates an image given a prompt.
func MakeCreationRequest(request *CreationRequest, organizationID *string) (*Response, error) {
	return common.MakeRequest[CreationRequest, Response](request, CreateEndpoint, http.MethodPost, organizationID)
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

	prompt, err := writer.CreateFormField("prompt")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(prompt, strings.NewReader(request.Prompt))
	if err != nil {
		return nil, err
	}

	n, err := writer.CreateFormField("n")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(n, strings.NewReader(fmt.Sprintf("%d", request.N)))
	if err != nil {
		return nil, err
	}

	size, err := writer.CreateFormField("size")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(size, strings.NewReader(request.Size))
	if err != nil {
		return nil, err
	}

	respFmt, err := writer.CreateFormField("response_format")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(respFmt, strings.NewReader(request.ResponseFormat))
	if err != nil {
		return nil, err
	}

	user, err := writer.CreateFormField("user")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(user, strings.NewReader(request.User))
	if err != nil {
		return nil, err
	}

	err = common.CreateFormFile("image", request.ImageName, request.Image, writer)
	if err != nil {
		return nil, err
	}

	err = common.CreateFormFile("mask", request.MaskName, request.Mask, writer)
	if err != nil {
		return nil, err
	}

	writer.Close()
	return common.MakeRequestWithForm[Response](buf, EditEndpoint, http.MethodPost, writer.FormDataContentType(), organizationID)
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

	n, err := writer.CreateFormField("n")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(n, strings.NewReader(fmt.Sprintf("%d", request.N)))
	if err != nil {
		return nil, err
	}

	size, err := writer.CreateFormField("size")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(size, strings.NewReader(request.Size))
	if err != nil {
		return nil, err
	}

	respFmt, err := writer.CreateFormField("response_format")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(respFmt, strings.NewReader(request.ResponseFormat))
	if err != nil {
		return nil, err
	}

	user, err := writer.CreateFormField("user")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(user, strings.NewReader(request.User))
	if err != nil {
		return nil, err
	}

	err = common.CreateFormFile("image", request.ImageName, request.Image, writer)
	if err != nil {
		return nil, err
	}

	writer.Close()
	return common.MakeRequestWithForm[Response](buf, VariationEndpoint, http.MethodPost, writer.FormDataContentType(), organizationID)
}
