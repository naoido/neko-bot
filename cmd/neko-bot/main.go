package main

import (
	"context"
	"fmt"
	"neko-bot/internal/bot"
	"neko-bot/internal/cli"
	"neko-bot/internal/client/ipa"
)

func main() {
	/*
		Select a stage to run the bot.
		stages:
			- prod: 本番環境
			- dev: 開発環境(ローカルで使用する場合はこちらを選択)
	*/
	fmt.Println("Create session of NEKO BOT.")

	loadingDone := make(chan struct{})
	cli.Loading(loadingDone)

	botClient := bot.NewBot()
	err := botClient.Start()
	if err != nil {
		panic(err)
		return
	}

	close(loadingDone)

	fmt.Println("\u001b[00;32m・▶ ︎Bot is now running.・\u001b[00m")
	fmt.Println("\u001B[00;31m・> Press q to exit.・\u001B[00m")

	ipa.StartWatch(botClient.Session())

	c := cli.NewCLI()
	ctx := context.Background()

	err = c.Start(ctx)
	if err != nil {
		panic(err)
	}
	<-c.Done()

	loadingDone = make(chan struct{})
	cli.Loading(loadingDone)

	err = botClient.Stop()
	if err != nil {
		panic(err)
		return
	}

	select {
	case <-botClient.Done():
		fmt.Println("\u001B[00;31mFinished.\u001B[00m")
	}
	close(loadingDone)
}
