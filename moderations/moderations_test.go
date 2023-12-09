package moderations_test

import (
	"os"
	"testing"

	"github.com/Kardbord/gopenai/authentication"
	"github.com/Kardbord/gopenai/moderations"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func TestModerations(t *testing.T) {
	resp, err := moderations.MakeRequest(&moderations.Request{
		Input: []string{"So long, and thanks for all the fish."},
		Model: moderations.ModelStable,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response received")
	}
}
