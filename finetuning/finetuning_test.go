package finetuning_test

import (
	"os"
	"testing"

	"github.com/Kardbord/gopenai/authentication"
	"github.com/Kardbord/gopenai/finetuning"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func TestFinetunes(t *testing.T) {
	// TODO: build more comprehensive tests for the finetunes endpoint.
	limit := uint64(2)
	after := "ft-event-TjX0lMfOniCZX64t9PUQT5hn"

	{ // Simple request
		_, err := finetuning.MakeListRequest(nil, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
	}
	{ // After
		_, err := finetuning.MakeListRequest(nil, &after, nil)
		if err != nil {
			t.Fatal(err)
		}
	}
	{ // Limit
		_, err := finetuning.MakeListRequest(&limit, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
	}
	{ // After and Limit
		_, err := finetuning.MakeListRequest(&limit, &after, nil)
		if err != nil {
			t.Fatal(err)
		}
	}
}
