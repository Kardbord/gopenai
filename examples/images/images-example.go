package main

import (
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

	fmt.Printf("Creating from prompt: %s\n", prompt)
	resp, _, err := images.MakeModeratedCreationRequest(&images.CreationRequest{
		Prompt: prompt,
		Size:   images.Dalle2SmallImage,
		User:   "https://github.com/TannerKvarfordt/gopenai",
	}, nil)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Generated: %s\n", resp.Data[0].URL)
	return resp, nil
}

func variation(imagename, image string) error {

	fmt.Printf("Generating a variation...")
	resp, err := images.MakeVariationRequest(&images.VariationRequest{
		Image:     image,
		ImageName: imagename,
		Size:      images.Dalle2SmallImage,
		User:      "https://github.com/TannerKvarfordt/gopenai",
	}, nil)
	if err != nil {
		return err
	}

	fmt.Printf("Generated: %s\n", resp.Data[0].URL)
	return nil
}

func main() {
	resp, err := create()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = variation("Original", resp.Data[0].URL)
	if err != nil {
		fmt.Println(err)
		return
	}
}
