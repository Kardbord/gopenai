package images_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/TannerKvarfordt/gopenai/authentication"
	"github.com/TannerKvarfordt/gopenai/images"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func create(model, size string) (*images.Response, error) {
	const prompt = "A cute baby sea otter"

	fmt.Printf("Creating from prompt: %s\n", prompt)
	resp, err := images.MakeCreationRequest(&images.CreationRequest{
		Prompt: prompt,
		Size:   size,
		User:   "https://github.com/TannerKvarfordt/gopenai",
		Model:  model,
	}, nil)
	if err != nil {
		return nil, err
	}
	if len(resp.Data) < 1 {
		return nil, errors.New("no images created")
	}

	fmt.Printf("Generated: %s\n", resp.Data[0].URL)
	return resp, nil
}

func variation(model, imagename, image string) error {

	fmt.Printf("Generating a variation...")
	resp, err := images.MakeVariationRequest(&images.VariationRequest{
		Image:     image,
		ImageName: imagename,
		Size:      images.Dalle2SmallImage,
		User:      "https://github.com/TannerKvarfordt/gopenai",
		Model:     model,
	}, nil)
	if err != nil {
		return err
	}
	if len(resp.Data) < 1 {
		return errors.New("no images edited")
	}

	fmt.Printf("Generated: %s\n", resp.Data[0].URL)
	return nil
}

func TestImagesDalle2(t *testing.T) {
	resp, err := create(images.ModelDalle2, images.Dalle2SmallImage)
	if err != nil {
		t.Fatal(err)
	}

	err = variation(images.ModelDalle2, "Original", resp.Data[0].URL)
	if err != nil {
		t.Fatal(err)
	}
}

func TestImagesDalle3(t *testing.T) {
	_, err := create(images.ModelDalle3, images.Dalle3SquareImage)
	if err != nil {
		t.Fatal(err)
	}
}
