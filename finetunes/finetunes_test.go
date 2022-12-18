package finetunes_test

import (
	"os"
	"testing"

	"github.com/TannerKvarfordt/gopenai/authentication"
	"github.com/TannerKvarfordt/gopenai/finetunes"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func TestFinetunes(t *testing.T) {
	// TODO: build a more comprehensive tests for the finetunes endpoint.

	resp, err := finetunes.MakeListRequest(nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response received")
	}
	if resp.Error != nil {
		t.Fatalf("%s -> %s\n", resp.Error.Type, resp.Error.Message)
	}
}
