package finetunes_test

import (
	"os"
	"testing"

	"github.com/Kardbord/gopenai/authentication"
	"github.com/Kardbord/gopenai/finetunes"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func TestFinetunes(t *testing.T) {
	// TODO: build a more comprehensive tests for the finetunes endpoint.

	_, err := finetunes.MakeListRequest(nil)
	if err != nil {
		t.Fatal(err)
	}
}
