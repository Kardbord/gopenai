package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/TannerKvarfordt/gopenai/authentication"
	"github.com/TannerKvarfordt/gopenai/files"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"
const file = "testdata.jsonl"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func list() error {
	resp, err := files.MakeListRequest(nil)
	if err != nil {
		return err
	}
	if resp == nil {
		return errors.New("nil response received")
	}
	if resp.Error != nil {
		return fmt.Errorf("%s -> %s", resp.Error.Type, resp.Error.Message)
	}

	fmt.Printf("You currently have %d files uploaded to Open AI.\n", len(resp.Data))
	return nil
}

func upload() (string, error) {
	resp, err := files.MakeUploadRequest(&files.UploadRequest{
		Purpose:  "fine-tune",
		Filename: file,
		Filepath: "./" + file,
	}, nil)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("nil response received")
	}
	if resp.Error != nil {
		return "", fmt.Errorf("%s -> %s", resp.Error.Type, resp.Error.Message)
	}

	fmt.Printf("Uploaded %s to Open AI, with purpose=\"%s\"\n", resp.Filename, resp.Purpose)
	return resp.ID, nil
}

func retrieve(fileID string) error {
	resp, err := files.MakeRetrieveRequest(fileID, nil)
	if err != nil {
		return err
	}
	if resp == nil {
		return errors.New("nil response received")
	}
	if resp.Error != nil {
		return fmt.Errorf("%s -> %s", resp.Error.Type, resp.Error.Message)
	}

	fmt.Printf("Retrieved fileID=%s from Open AI, with name=\"%s\" purpose=\"%s\"\n", fileID, resp.Filename, resp.Purpose)
	return nil
}

func retrieveContent(fileID string) error {
	err := files.MakeRetrieveContentRequest(fileID, file, true, nil)
	if err != nil {
		return err
	}

	fmt.Printf("Retrieved fileID=%s content and wrote to disk at %s\n", fileID, file)
	return nil
}

func delete(fileID string) error {
	resp, err := files.MakeDeleteRequest(fileID, nil)
	if err != nil {
		return err
	}
	if resp == nil {
		return errors.New("nil response received")
	}
	if resp.Error != nil {
		return fmt.Errorf("%s -> %s", resp.Error.Type, resp.Error.Message)
	}
	if !resp.Deleted {
		return errors.New("failed to delete remote file")
	}

	fmt.Printf("Deleted %s from the Open AI servers.\n", fileID)
	return nil
}

func main() {
	var err error

	err = list()
	if err != nil {
		fmt.Println(err)
		return
	}

	var fileID string
	fileID, err = upload()
	if err != nil {
		fmt.Println(err)
		return
	}

	const sleepDuration = 5
	for i := 0; i < sleepDuration; i++ {
		fmt.Printf("Sleeping to allow the file to process %d/%ds\n", i, sleepDuration)
		time.Sleep(time.Second)
	}

	err = list()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = retrieve(fileID)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = retrieveContent(fileID)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = delete(fileID)
	if err != nil {
		fmt.Println(err)
		return
	}
}
