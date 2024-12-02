package main

import (
	"fmt"
	"os"

	"github.com/Kardbord/gopenai/authentication"
	"github.com/Kardbord/gopenai/completions"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func main() {
	prompt := "So long, and thanks for all the"
	resp, _, err := completions.MakeModeratedRequest(&completions.Request{
		Model:            "gpt-3.5-turbo-instruct",
		Prompt:           []string{prompt},
		MaxTokens:        5,
		Echo:             true,
		PresencePenalty:  0.0,
		FrequencyPenalty: 0.0,
		BestOf:           nil,
		User:             "https://github.com/Kardbord/gopenai",
	}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Prompt: %s\nResponse: %s\n", prompt, resp.Choices[0].Text)
}
