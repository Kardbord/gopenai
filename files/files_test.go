package files_test

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/TannerKvarfordt/gopenai/authentication"
	"github.com/TannerKvarfordt/gopenai/files"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"
const file = "./test_files/testdata.jsonl"
const file2 = "retrieveddata2.jsonl"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func list(t *testing.T) error {
	resp, err := files.MakeListRequest(nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal(errors.New("nil response received"))
	}
	if resp.Error != nil {
		t.Fatal(fmt.Errorf("%s -> %s", resp.Error.Type, resp.Error.Message))
	}
	return nil
}

func upload(t *testing.T) (string, error) {
	resp, err := files.MakeUploadRequest(&files.UploadRequest{
		Purpose:  "fine-tune",
		Filename: file,
		Filepath: "./" + file,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal(errors.New("nil response received"))
	}
	if resp.Error != nil {
		t.Fatal(fmt.Errorf("%s -> %s", resp.Error.Type, resp.Error.Message))
	}

	return resp.ID, nil
}

func retrieve(fileID string, t *testing.T) error {
	resp, err := files.MakeRetrieveRequest(fileID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal(errors.New("nil response received"))
	}
	if resp.Error != nil {
		t.Fatal(fmt.Errorf("%s -> %s", resp.Error.Type, resp.Error.Message))
	}

	return nil
}

func retrieveContent(fileID string, t *testing.T) error {
	err := files.MakeRetrieveContentRequest(fileID, file2, true, nil)
	if err != nil {
		t.Fatal(err)
	}

	return nil
}

func delete(fileID string, t *testing.T) error {
	resp, err := files.MakeDeleteRequest(fileID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal(errors.New("nil response received"))
	}
	if resp.Error != nil {
		t.Fatal(fmt.Errorf("%s -> %s", resp.Error.Type, resp.Error.Message))
	}
	if !resp.Deleted {
		t.Fatal(errors.New("failed to delete remote file"))
	}

	_ = os.Remove(file2)
	return nil
}

func TestFiles(t *testing.T) {
	var err error

	err = list(t)
	if err != nil {
		fmt.Println(err)
		return
	}

	var fileID string
	fileID, err = upload(t)
	if err != nil {
		fmt.Println(err)
		return
	}

	const sleepDuration = 5
	for i := 0; i < sleepDuration; i++ {
		fmt.Printf("Sleeping to allow the file to process %d/%ds\n", i, sleepDuration)
		time.Sleep(time.Second)
	}

	err = list(t)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = retrieve(fileID, t)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = retrieveContent(fileID, t)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = delete(fileID, t)
	if err != nil {
		fmt.Println(err)
		return
	}
}
