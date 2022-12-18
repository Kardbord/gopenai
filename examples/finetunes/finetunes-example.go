package main

import (
	"fmt"
	"os"

	"github.com/TannerKvarfordt/gopenai/authentication"
	"github.com/TannerKvarfordt/gopenai/finetunes"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func main() {
	// TODO: build a more comprehensive example of how to use this endpoint.

	resp, err := finetunes.MakeListRequest(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp == nil {
		fmt.Println("nil response received")
		return
	}
	if resp.Error != nil {
		fmt.Printf("%s -> %s\n", resp.Error.Type, resp.Error.Message)
		return
	}

	fmt.Printf("You currently have %d finetune jobs running.\n", len(resp.Data))
}
