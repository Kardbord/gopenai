package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/TannerKvarfordt/gopenai/authentication"
	"github.com/TannerKvarfordt/gopenai/images"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func create() (*images.Response, error) {
	const prompt = "A cute baby sea otter"

	fmt.Printf("Prompt: %s\n", prompt)
	resp, err := images.MakeCreationRequest(&images.CreationRequest{
		Prompt: prompt,
		Size:   images.SmallImage,
		User:   "https://github.com/TannerKvarfordt/gopenai",
	}, nil)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("%s -> %s", resp.Error.Type, resp.Error.Message)
	}
	if len(resp.Data) < 1 {
		return nil, errors.New("no images created")
	}

	fmt.Printf("Generated: %s\n", resp.Data[0].URL)
	return resp, nil
}

func main() {
	_, err := create()
	if err != nil {
		fmt.Println(err)
		return
	}
}
