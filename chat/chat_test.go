package chat_test

import (
	"os"
	"testing"

	"github.com/TannerKvarfordt/gopenai/authentication"
	"github.com/TannerKvarfordt/gopenai/chat"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

func TestChat(t *testing.T) {
	resp, err := chat.MakeRequest(&chat.Request{
		Model: "gpt-3.5-turbo",
		Messages: []chat.Chat{
			{
				Role:    "user",
				Content: "Hello!",
			},
		},
		User: "https://github.com/TannerKvarfordt/gopenai",
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
