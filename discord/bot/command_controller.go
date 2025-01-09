package bot

import (
	"neko-bot/discord/command"
	"neko-bot/discord/handler"
)

func RegisterCommands() error {
	// ここに使用するコマンドを登録していく
	handler.Add(&command.NewPing("ping", &config.Prefix).Detail)
	handler.Add(&command.NewVoice("voice", &config.Prefix).Detail)
	handler.Add(&command.NewSetting("setting", &config.Prefix).Detail)
	handler.Add(&command.NewSay("say", &config.Prefix).Detail)

	return handler.RegisterCommands(bot.Session())
}

func RemoveCommands() error {
	return handler.RemoveCommands(bot.Session())
}
