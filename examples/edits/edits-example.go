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
	resp, err := edits.MakeRequest(&edits.Request{
		Model:       "text-davinci-edit-001",
		Input:       input,
		Instruction: "Fix the spelling mistakes",
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
	if len(resp.Choices) < 1 {
		fmt.Println("No choices received")
		return
	}

	fmt.Printf("Input=%s\n", input)
	fmt.Printf("Result=%s\n", resp.Choices[0].Text)
}
