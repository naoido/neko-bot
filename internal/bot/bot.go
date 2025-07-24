package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"neko-bot/internal/bot/handler"
	"neko-bot/internal/config"
	"time"
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

type Bot interface {
	Session() *discordgo.Session
	StartedAt() time.Time
	Start() error
	Stop() error
	Done() <-chan struct{}
}

type bot struct {
	session            *discordgo.Session
	startedAt          time.Time
	registeredCommands []*discordgo.ApplicationCommand
	done               chan struct{}
}

func NewBot() Bot {
	session, err := discordgo.New(config.BotConfig().Token())
	if err != nil {
		panic(err)
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsAll | discordgo.PermissionSendMessages

	return &bot{
		session:   session,
		startedAt: time.Now(),
		done:      make(chan struct{}),
	}
}

func (b *bot) Session() *discordgo.Session {
	return b.session
}

func (b *bot) StartedAt() time.Time {
	return b.startedAt
}

func (b *bot) Start() error {
	fmt.Print(nekoBot)

	hdl := handler.NewHandler(b.session)

	// ハンドラーの追加
	for _, h := range hdl.Handlers() {
		b.session.AddHandler(h)
	}

	err := b.session.Open()
	if err != nil {
		return err
	}

	// コマンドの登録
	for _, model := range hdl.Commands() {
		cmd, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, model.GetCommandData().GuildID, model.GetCommandData())
		if err != nil {
			return err
		}
		b.registeredCommands = append(b.registeredCommands, cmd)
		fmt.Printf("\rRegistered command %s\r\n", cmd.Name)
	}

	return nil
}

func (b *bot) Stop() error {
	defer close(b.done)

	fmt.Println("\r\u001B[00;31mShutting down...\u001B[00m")
	for _, model := range b.registeredCommands {
		err := b.session.ApplicationCommandDelete(b.session.State.User.ID, model.GuildID, model.ID)
		if err != nil {
			return err
		}
		fmt.Printf("\rRemove command %s\r\n", model.Name)
	}

	return b.session.Close()
}

func (b *bot) Done() <-chan struct{} {
	return b.done
}
