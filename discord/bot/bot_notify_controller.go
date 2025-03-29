package bot

import (
	"neko-bot/discord/handler"
)

func RegisterHandlers() {
	handler.RegisterHandlers(bot.Session())
}
