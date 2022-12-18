package edits_test

import (
	"os"
	"testing"

	"github.com/TannerKvarfordt/gopenai/authentication"
	"github.com/TannerKvarfordt/gopenai/edits"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func TestEdits(t *testing.T) {
	resp, err := edits.MakeRequest(&edits.Request{
		Model:       "text-davinci-edit-001",
		Input:       "What day of the wek is it?",
		Instruction: "Fix the spelling mistakes",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Choices) < 1 {
		t.Fatal("No choices received")
	}
}
