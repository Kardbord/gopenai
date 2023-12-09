package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Kardbord/gopenai/authentication"
	"github.com/Kardbord/gopenai/chat"
)

const OpenAITokenEnv = "OPENAI_API_KEY"

func init() {
	key := os.Getenv(OpenAITokenEnv)
	authentication.SetAPIKey(key)
}

const (
	model string = "gpt-3.5-turbo"
	quit  string = "quit"
)

func main() {
	fmt.Println("Welcome to ChatGPT! Quit at any time by typing \"quit\".\nStart the conversation by typing a message, and hitting enter.")

	chatlog := make([]chat.Chat, 0, 10)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("User: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if strings.ToLower(strings.TrimSpace(input)) == quit {
			fmt.Println("Goodbye!")
			return
		}

		fmt.Println("ChatGPT is thinking...")
		r, _, err := chat.MakeModeratedRequest(&chat.Request{
			Model: model,
			Messages: append(chatlog, chat.Chat{
				Role:    chat.UserRole,
				Content: input,
			}),
			User: "https://github.com/Kardbord/gopenai",
		}, nil)

		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

		chatlog = append(chatlog, r.Choices[0].Message)
		fmt.Println("ChatGPT:", r.Choices[0].Message.Content)
	}
}
