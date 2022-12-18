// Package common contains common API structures and helper functions, not specific to an endpoint or model.
package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	auth "github.com/TannerKvarfordt/gopenai/authentication"
)

const (
	// The version of the API currently implemented by this library.
	APIVersion = "v1"

	// The basis of all API endpoints.
	BaseURL = "https://api.openai.com/" + APIVersion + "/"
)

// A common error structure included in OpenAI API response bodies.
type ResponseError struct {
	// The error message.
	Message string `json:"message"`
	// The error type.
	Type string `json:"type"`
}

// A common usage information structure included in OpenAI API response bodies.
type ResponseUsage struct {
	PromptTokens     uint64 `json:"prompt_tokens"`
	CompletionTokens uint64 `json:"completion_tokens"`
	TotalTokens      uint64 `json:"total_tokens"`
}

// Send a request to the given OpenAI endpoint.
// The method parameter should be an HTTP method, such as GET or POST.
// The organizationID parameter is optional. If provided, it will be included in the request header.
// If not provided, the authorization.DefaultOrganizationID will be used, if it is set.
func MakeRequest[RequestT any, ResponseT any](request *RequestT, endpoint, method string, organizationID *string) (*ResponseT, error) {
	var req *http.Request = nil
	var err error = nil
	if request != nil {
		jsonData, err2 := json.Marshal(request)
		if err2 != nil {
			return nil, err2
		}
		req, err = http.NewRequest(method, endpoint, bytes.NewBuffer(jsonData))
	} else {
		req, err = http.NewRequest(method, endpoint, nil)
	}
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errors.New("nil request created")
	}
	SetRequestHeaders(req, "application/json", organizationID)
	return makeRequest[ResponseT](req)
}

func MakeRequestWithForm[ResponseT any](form *bytes.Buffer, endpoint, method, contentType string, organizationID *string) (*ResponseT, error) {
	req, err := http.NewRequest(method, endpoint, form)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errors.New("nil request created")
	}

	SetRequestHeaders(req, contentType, organizationID)
	return makeRequest[ResponseT](req)
}

func SetRequestHeaders(req *http.Request, contentType string, organizationID *string) {
	if req == nil {
		return
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set(auth.AuthHeaderKey, auth.AuthHeaderPrefix+auth.APIKey())

	if organizationID != nil {
		req.Header.Set(auth.OrgHeaderKey, *organizationID)
	} else if len(auth.DefaultOrganizationID()) != 0 {
		req.Header.Set(auth.OrgHeaderKey, auth.DefaultOrganizationID())
	}
}

func makeRequest[ResponseT any](req *http.Request) (*ResponseT, error) {
	if req == nil {
		return nil, errors.New("nil request provided to makeRequest helper - this is a bug in the library")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("nil response received")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if respBody == nil {
		return nil, errors.New("unable to parse response body")
	}

	response := new(ResponseT)
	err = json.Unmarshal(respBody, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func CreateFormFile(fieldname, filename, filepath string, writer *multipart.Writer) error {
	file, err := writer.CreateFormFile(fieldname, filename)
	if err != nil {
		return err
	}

	var fdata io.ReadCloser
	if IsUrl(filepath) {
		resp, err := http.Get(filepath)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to retrieve file from url, status code: %d", resp.StatusCode)
		}
		fdata = resp.Body
	} else {
		fdata, err = os.Open(filepath)
		if err != nil {
			return err
		}
	}
	defer fdata.Close()

	_, err = io.Copy(file, fdata)
	if err != nil {
		return err
	}

	return nil
}
