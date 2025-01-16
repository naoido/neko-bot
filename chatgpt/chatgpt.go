package chatgpt

import (
	"context"
	"github.com/ayush6624/go-chatgpt"
	"neko-bot/discord/neko"
)

var ctx = context.Background()
var client *chatgpt.Client

func init() {
	key := neko.GetConfig().ChatgptKey

	c, err := chatgpt.NewClient(key)
	if err != nil {
		panic(err)
	}

	client = c
}

func GetResponse(message string) (*chatgpt.ChatResponse, error) {
	return client.Send(ctx, &chatgpt.ChatCompletionRequest{
		Model: chatgpt.GPT4,
		Messages: []chatgpt.ChatMessage{
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: "語尾に必ず「なのだ」をつけてください",
			},
			{
				Role:    chatgpt.ChatGPTModelRoleUser,
				Content: message,
			},
		},
	})
}
