package embeddings_test

import (
	"os"
	"testing"

	"github.com/Kardbord/gopenai/authentication"
	"github.com/Kardbord/gopenai/embeddings"
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
		User:  "https://github.com/Kardbord/gopenai",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Data[0].Embedding) < 1 {
		t.Fatal("Empty embedding returned in response")
	}
}
