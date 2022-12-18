package main

import (
	"fmt"
	"os"

	"github.com/TannerKvarfordt/gopenai/authentication"
	"github.com/TannerKvarfordt/gopenai/completions"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func main() {
	prompt := "So long, and thanks for all the"
	resp, err := completions.MakeRequest(&completions.Request{
		Model:            "text-ada-001",
		Prompt:           []string{prompt},
		MaxTokens:        5,
		Echo:             true,
		PresencePenalty:  0.0,
		FrequencyPenalty: 0.0,
		BestOf:           nil,
		User:             "https://github.com/TannerKvarfordt/gopenai",
	}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Prompt: %s\nResponse: %s\n", prompt, resp.Choices[0].Text)
}
