package embeddings_test

import (
	"os"
	"testing"

	"github.com/TannerKvarfordt/gopenai/authentication"
	"github.com/TannerKvarfordt/gopenai/embeddings"
	_ "github.com/joho/godotenv/autoload"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func TestEmbeddings(t *testing.T) {
	input := "The food was delicious and the waiter..."
	resp, err := embeddings.MakeRequest(&embeddings.Request{
		Model: "text-embedding-ada-002",
		Input: []string{input},
		User:  "https://github.com/TannerKvarfordt/gopenai",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response received")
	}
	if resp.Error != nil {
		t.Fatalf("%s -> %s\n", resp.Error.Type, resp.Error.Message)
	}
	if len(resp.Data) < 1 {
		t.Fatal("No data returned in response")
	}
	if len(resp.Data[0].Embedding) < 1 {
		t.Fatal("Empty embedding returned in response")
	}
}
