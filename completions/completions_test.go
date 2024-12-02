package completions_test

import (
	"os"
	"testing"

	"github.com/Kardbord/gopenai/authentication"
	"github.com/Kardbord/gopenai/completions"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func TestCompletions(t *testing.T) {
	resp, err := completions.MakeRequest(&completions.Request{
		Model:            "gpt-3.5-turbo-instruct",
		Prompt:           []string{"So long, and thanks for all the"},
		MaxTokens:        5,
		Echo:             false,
		PresencePenalty:  0.0,
		FrequencyPenalty: 0.0,
		User:             "https://github.com/Kardbord/gopenai",
	}, nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	if len(resp.Choices) < 1 {
		t.Fatal("no choices received")
		return
	}
}
