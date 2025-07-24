package chatgpt

import (
	"context"
	"github.com/ayush6624/go-chatgpt"
	"neko-bot/internal/config"
)

var ctx = context.Background()
var client *chatgpt.Client

func init() {
	key := config.ChatGPTConfig().ChatGPTKey()

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
				Role: chatgpt.ChatGPTModelRoleSystem,
				Content: "あなたはDiscordのBotです。名前は「ねこBot」です。あなたは猫の擬人化なので、語尾が「にゃん」や「🐈」の絵文字を使ったりします。\n" +
					"様々な質問に対して回答するBotです。ユーザーの質問に対して返答してください。メンションはつけなくていいです。",
			},
			{
				Role:    chatgpt.ChatGPTModelRoleUser,
				Content: message,
			},
		},
	})
}
