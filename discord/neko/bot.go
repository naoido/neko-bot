package neko

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"neko-bot/internal/errors"
	"neko-bot/internal/loading"
)

const nekoBot = `
          _____                    _____                    _____                   _______                   _____                   _______               _____
         /\    \                  /\    \                  /\    \                 /::\    \                 /\    \                 /::\    \             /\    \
        /::\____\                /::\    \                /::\____\               /::::\    \               /::\    \               /::::\    \           /::\    \
       /::::|   |               /::::\    \              /:::/    /              /::::::\    \             /::::\    \             /::::::\    \          \:::\    \
      /:::::|   |              /::::::\    \            /:::/    /              /::::::::\    \           /::::::\    \           /::::::::\    \          \:::\    \
     /::::::|   |             /:::/\:::\    \          /:::/    /              /:::/~~\:::\    \         /:::/\:::\    \         /:::/~~\:::\    \          \:::\    \
    /:::/|::|   |            /:::/__\:::\    \        /:::/____/              /:::/    \:::\    \       /:::/__\:::\    \       /:::/    \:::\    \          \:::\    \
   /:::/ |::|   |           /::::\   \:::\    \      /::::\    \             /:::/    / \:::\    \     /::::\   \:::\    \     /:::/    / \:::\    \         /::::\    \
  /:::/  |::|   | _____    /::::::\   \:::\    \    /::::::\____\________   /:::/____/   \:::\____\   /::::::\   \:::\    \   /:::/____/   \:::\____\       /::::::\    \
 /:::/   |::|   |/\    \  /:::/\:::\   \:::\    \  /:::/\:::::::::::\    \ |:::|    |     |:::|    | /:::/\:::\   \:::\ ___\ |:::|    |     |:::|    |     /:::/\:::\    \
/:: /    |::|   /::\____\/:::/__\:::\   \:::\____\/:::/  |:::::::::::\____\|:::|____|     |:::|    |/:::/__\:::\   \:::|    ||:::|____|     |:::|    |    /:::/  \:::\____\
\::/    /|::|  /:::/    /\:::\   \:::\   \::/    /\::/   |::|~~~|~~~~~      \:::\    \   /:::/    / \:::\   \:::\  /:::|____| \:::\    \   /:::/    /    /:::/    \::/    /
 \/____/ |::| /:::/    /  \:::\   \:::\   \/____/  \/____|::|   |            \:::\    \ /:::/    /   \:::\   \:::\/:::/    /   \:::\    \ /:::/    /    /:::/    / \/____/
         |::|/:::/    /    \:::\   \:::\    \            |::|   |             \:::\    /:::/    /     \:::\   \::::::/    /     \:::\    /:::/    /    /:::/    /
         |::::::/    /      \:::\   \:::\____\           |::|   |              \:::\__/:::/    /       \:::\   \::::/    /       \:::\__/:::/    /    /:::/    /
         |:::::/    /        \:::\   \::/    /           |::|   |               \::::::::/    /         \:::\  /:::/    /         \::::::::/    /     \::/    /
         |::::/    /          \:::\   \/____/            |::|   |                \::::::/    /           \:::\/:::/    /           \::::::/    /       \/____/
         /:::/    /            \:::\    \                |::|   |                 \::::/    /             \::::::/    /             \::::/    /
        /:::/    /              \:::\____\               \::|   |                  \::/____/               \::::/    /               \::/____/
        \::/    /                \::/    /                \:|   |                   ~~                      \::/____/                 ~~
         \/____/                  \/____/                  \|___|                                            ~~

`

type Bot struct {
	discord *discordgo.Session
}

func New(config Config) (*Bot, error) {
	bot, err := start(config)
	errors.CatchAndPanic(err, "cannot open discord bot session")

	return &Bot{
		discord: bot.Session(),
	}, nil
}

func start(config Config) (*Bot, error) {
	// Create discord session
	session, err := discordgo.New(config.token)
	if err != nil {
		return nil, err
	}

	// Open discord session
	err = session.Open()
	if err != nil {
		return nil, err
	}

	b := &Bot{discord: session}

	b.UpdateBot(config, false)
	fmt.Print(nekoBot)

	return b, nil
}

func (b Bot) Session() *discordgo.Session {
	if b.discord == nil {
		panic("Discord not initialized")
	}
	return b.discord
}

func (b Bot) Stop() error {
	fmt.Println("\r\u001B[00;31mShutting down...\u001B[00m")
	loading.Start()
	err := b.discord.Close()
	loading.Stop()
	return err
}

func (b Bot) UpdateBot(config Config, reload bool) {
	if reload {
		var err error
		errors.Catch(err, "could not get config")

		err = b.Session().Close()
		errors.Catch(err, "could not close discord")

		err = b.Session().Open()
		errors.CatchAndPanic(err, "could not open discord")
	}
	err := b.discord.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: config.status,
		Activities: []*discordgo.Activity{
			{
				Name: config.activeMessage,
				Type: discordgo.ActivityTypeGame,
			},
		},
	})
	errors.Catch(err, "\rfailed to update status of discord bot\r")
}
