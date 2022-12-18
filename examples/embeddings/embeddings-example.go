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
	resp, _, err := embeddings.MakeModeratedRequest(&embeddings.Request{
		Model: "text-embedding-ada-002",
		Input: []string{input},
		User:  "https://github.com/TannerKvarfordt/gopenai",
	}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Input: %s\n", input)
	fmt.Printf("Resulted in embedding size: %d\n", len(resp.Data[0].Embedding))
}
