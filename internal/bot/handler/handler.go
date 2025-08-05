package handler

import (
	"github.com/bwmarrin/discordgo"
	"neko-bot/internal/bot/command"
)

type Handler interface {
	Handlers() []interface{}
	Commands() []command.Command
}

type handler struct{}

var (
	handlers []interface{}
	commands []command.Command
)

func NewHandler(s *discordgo.Session) Handler {
	// コマンド一覧
	commands = []command.Command{
		command.NewLastAction(),
		command.NewPing(),
		command.NewSay(),
		command.NewSetting(),
		command.NewVoice(),
		command.NewZenn(s),
	}

	// イベントリスナー一覧
	handlers = []interface{}{
		onCommand,
		onMessage,
		onThreadCreate,
		onReaction,
		onReactionRemove,
	}

	return &handler{}
}

func (h *handler) Handlers() []interface{} {
	return handlers
}

func (h *handler) Commands() []command.Command {
	return commands
}
