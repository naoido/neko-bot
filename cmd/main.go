package main

import (
	"fmt"
	"neko-bot/internal/listening"
	"neko-bot/internal/zr"
	"neko-bot/neko"
	"os"
)

func main() {
	/*
		Select a stage to run the bot.
		stages:
			- prod: 本番環境
			- dev: 開発環境(ローカルで使用する場合はこちらを選択)
	*/
	fmt.Println("Create session of NEKO BOT.")
	stage := zr.OrDef(os.Getenv("STAGE"), "dev")
	neko.Start(stage)

	fmt.Println("\u001b[00;32m・▶ ︎Bot is now running.・\u001b[00m")
	fmt.Println("\u001B[00;31m・> Press q to exit.・\u001B[00m")

	listening.KeyListener()
	neko.Stop()
}
