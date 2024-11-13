package bot

import (
	"neko-bot/discord/command"
	"neko-bot/discord/handler"
)

func RegisterCommands() error {
	// ここに使用するコマンドを登録していく
	handler.Add(&command.NewPing("ping", &config.Prefix).Model)
	handler.Add(&command.NewUmbrella("umbrella", &config.Prefix).Model) 

	return handler.RegisterCommands(bot.Session())
}

func RemoveCommands() error {
	return handler.RemoveCommands(bot.Session())
}
