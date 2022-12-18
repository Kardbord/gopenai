// Package files provides bindings for the [files] [endpoint].
// Files are used to upload documents that can be used with
// features like Fine-tuning.
//
// [files]: https://beta.openai.com/docs/api-reference/files
// [endpoint]: https://api.openai.com/v1/files
package files

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/TannerKvarfordt/gopenai/common"
)

const Endpoint = common.BaseURL + "files"

type UploadedFile struct {
	ID        string                `json:"id"`
	Object    string                `json:"object"`
	Bytes     uint64                `json:"bytes"`
	CreatedAt uint64                `json:"created_at"`
	Filename  string                `json:"filename"`
	Purpose   string                `json:"purpose"`
	Error     *common.ResponseError `json:"error,omitempty"`
}

// Response structure for the files API "list" endpoint.
type ListResponse struct {
	Object string                `json:"object"`
	Data   []UploadedFile        `json:"data"`
	Error  *common.ResponseError `json:"error,omitempty"`
}

// Returns a list of files that belong to the user's organization.
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

// Request structure for the files "upload" endpoint.
type UploadRequest struct {
	// The intended purpose of the uploaded documents.
	// Use "fine-tune" for Fine-tuning. This allows us to validate the
	// format of the uploaded file.
	Purpose string `json:"purpose"`

	// Name of the [JSON Lines] file to be uploaded.
	//
	// If the purpose is set to "fine-tune", each line is a JSON record with
	// "prompt" and "completion" fields representing your training examples.
	//
	// Note that this is not the path to the file, but just the name of the file.
	//
	// [JSON Lines]: https://jsonlines.org/
	Filename string `json:"file"`

	// The path to the file, including the file's name and extension.
	Filepath string `json:"-"`
}

// Upload a file that contains document(s) to be used across various endpoints/features.
// Currently, the size of all the files uploaded by one organization can be up to 1 GB.
func MakeUploadRequest(request *UploadRequest, organizationID *string) (*UploadedFile, error) {
	// Implementation largely taken from https://github.com/sashabaranov/go-gpt3/blob/1c20931ead68f5d7e7e04747720fac1ebd73d35c/files.go#L53-L117

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	if len(request.Purpose) > 0 {
		err := common.CreateFormField("purpose", request.Purpose, writer)
		if err != nil {
			return nil, err
		}
	}

	if len(request.Filepath) > 0 {
		err := common.CreateFormFile("file", request.Filename, request.Filepath, writer)
		if err != nil {
			return nil, err
		}
	}

	writer.Close()
	r, err := common.MakeRequestWithForm[UploadedFile](buf, Endpoint, http.MethodPost, writer.FormDataContentType(), organizationID)
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

// Response structure for the files "delete" endpoint.
type DeleteResponse struct {
	ID      string                `json:"id"`
	Object  string                `json:"object"`
	Deleted bool                  `json:"deleted"`
	Error   *common.ResponseError `json:"error,omitempty"`
}

// Delete an uploaded file.
func MakeDeleteRequest(fileID string, organizationID *string) (*DeleteResponse, error) {
	r, err := common.MakeRequest[any, DeleteResponse](nil, fmt.Sprintf("%s/%s", Endpoint, fileID), http.MethodDelete, organizationID)
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

// Returns information about a specific file.
func MakeRetrieveRequest(fileID string, organizationID *string) (*UploadedFile, error) {
	r, err := common.MakeRequest[any, UploadedFile](nil, fmt.Sprintf("%s/%s", Endpoint, fileID), http.MethodGet, organizationID)
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

// Retreives "fileID" from Open AI, and writes it to disk at "filepath".
// If "filepath" already exists and "overwrite" is false, an error will be returned.
// If "filepath" already exists and "overwrite" is true, the existing file is truncated.
func MakeRetrieveContentRequest(fileID, filepath string, overwrite bool, organizationID *string) error {
	_, err := os.Stat(filepath)
	if err == nil && !overwrite {
		return os.ErrExist
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/content", Endpoint, fileID), nil)
	if err != nil {
		return err
	}
	if req == nil {
		return errors.New("nil request created")
	}

	common.SetRequestHeaders(req, "application/json", organizationID)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp == nil {
		return errors.New("nil response received")
	}
	defer resp.Body.Close()

	fout, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fout.Close()

	_, err = io.Copy(fout, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
