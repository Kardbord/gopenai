package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/TannerKvarfordt/gopenai/authentication"
	"github.com/TannerKvarfordt/gopenai/moderations"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func request(input string) {
	fmt.Println("Sending moderations request with input: ", input)

	resp, err := moderations.MakeRequest(&moderations.Request{
		Input: []string{input},
		Model: moderations.ModelStable,
	}, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	jsonResp, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		fmt.Printf(`Failed to marshal response to JSON for printing, err="%s"\n`, err)
		return
	}
	fmt.Printf("Response: %s\n", string(jsonResp))
}

func main() {
	request("I want to frolick in the sunshine with my friends.")
	request("I want to kill them.")
}
