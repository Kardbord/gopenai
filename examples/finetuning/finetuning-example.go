package main

import (
	"fmt"
	"os"

	"github.com/Kardbord/gopenai/authentication"
	"github.com/Kardbord/gopenai/finetuning"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func main() {
	// TODO: build a more comprehensive example of how to use this endpoint.

	resp, err := finetuning.MakeListRequest(nil, nil, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("You currently have %d finetune jobs running.\n", len(resp.Data))
}
