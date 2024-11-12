package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type Command interface {
	GetName() string
	GetCommandData() *discordgo.ApplicationCommand
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
	Prefix(s *discordgo.Session, m *discordgo.MessageCreate)
}

type Model struct {
	Command
	Detail Detail
}

type Detail struct {
	name   string
	prefix *string
}

func (d Detail) isCommand(i *discordgo.InteractionCreate) bool {
	return i.ApplicationCommandData().Name == d.name
}

func (d Detail) isPrefix(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	return m.Author.ID != s.State.User.ID && strings.HasPrefix(m.Content, fmt.Sprintf("%s%s", *d.prefix, d.name))
}
