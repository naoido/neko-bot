package command

import "github.com/bwmarrin/discordgo"

type Command interface {
	GetName() string
	Handler(s *discordgo.Session, m *discordgo.InteractionCreate)
	Prefix(s *discordgo.Session, i *discordgo.MessageCreate)
	Chat(s *discordgo.Session, i *discordgo.MessageCreate)
}
