package main

import (
	"fmt"
	"neko-bot/discord/bot"
	"neko-bot/internal/errors"
	"neko-bot/internal/listening"
	"neko-bot/internal/loading"
)

func main() {
	/*
		Select a stage to run the bot.
		stages:
			- prod: 本番環境
			- dev: 開発環境(ローカルで使用する場合はこちらを選択)
	*/
	fmt.Println("Create session of NEKO BOT.")

	loading.Start()
	err := bot.Start()
	errors.CatchAndPanic(err, "cannot start the bot")

	err = bot.Update()
	errors.CatchAndPanic(err, "cannot update the bot")
	loading.Stop()

	err = bot.RegisterCommands()
	errors.CatchAndPanic(err, "cannot register commands")

	fmt.Println("\u001b[00;32m・▶ ︎Bot is now running.・\u001b[00m")
	fmt.Println("\u001B[00;31m・> Press q to exit.・\u001B[00m")

	listening.KeyListener()
	err = bot.RemoveCommands()
	errors.CatchAndPanic(err, "cannot remove commands")

	err = bot.Stop()
	errors.CatchAndPanic(err, "cannot stop the bot")
}
