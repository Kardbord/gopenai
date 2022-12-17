package main

import (
	"fmt"
	"os"

	"github.com/TannerKvarfordt/gopenai/authentication"
	"github.com/TannerKvarfordt/gopenai/embeddings"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func main() {
	input := "The food was delicious and the waiter..."
	resp, err := embeddings.MakeRequest(&embeddings.Request{
		Model: "text-embedding-ada-002",
		Input: []string{input},
		User:  "https://github.com/TannerKvarfordt/gopenai",
	}, nil)
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
	if len(resp.Data) < 1 {
		fmt.Println("No data returned in response")
		return
	}

	fmt.Printf("Input: %s\n", input)
	fmt.Printf("Resulted in embedding size: %d\n", len(resp.Data[0].Embedding))
}
