package main

import (
	"fmt"
	"os"

	"github.com/TannerKvarfordt/gopenai/authentication"
	"github.com/TannerKvarfordt/gopenai/edits"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func main() {
	input := "What day of the wek is it?"
	resp, _, err := edits.MakeModeratedRequest(&edits.Request{
		Model:       "text-davinci-edit-001",
		Input:       input,
		Instruction: "Fix the spelling mistakes",
	}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Input=%s\n", input)
	fmt.Printf("Result=%s\n", resp.Choices[0].Text)
}
