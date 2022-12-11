// Package openaicommon contains common API structures and helper functions, not specific to an endpoint or model.
package openaicommon

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	auth "github.com/TannerKvarfordt/gopenai/authentication"
)

const (
	// The version of the API currently implemented by this library.
	APIVersion = "v1"

	// The basis of all API endpoints.
	BaseURL = "https://api.openai.com/" + APIVersion + "/"
)

func setRequestHeaders(req *http.Request, organizationID *string) {
	if req == nil {
		return
	}
	req.Header.Set(auth.AuthHeaderKey, auth.AuthHeaderPrefix+auth.APIKey())

	if organizationID != nil {
		req.Header.Set(auth.OrgHeaderKey, *organizationID)
	} else if len(auth.DefaultOrganizationID()) != 0 {
		req.Header.Set(auth.OrgHeaderKey, auth.DefaultOrganizationID())
	}
}

// Send a request containing the given JSON data to the given OpenAI endpoint.
// Optionally provide an organization ID to use. If not provided,
// authentication.DefaultOrganizationID() will be used if it is set.
func MakeRequest(jsonData []byte, endpoint, method string, organizationID *string) ([]byte, error) {
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errors.New("nil request created")
	}
	req.Header.Set("Content-Type", "application/json")
	setRequestHeaders(req, organizationID)

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

	return respBody, nil
}
