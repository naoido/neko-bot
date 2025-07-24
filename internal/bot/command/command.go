package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"neko-bot/internal/config"
	"strings"
)

type Command interface {
	GetName() string
	IsCommand(i *discordgo.InteractionCreate) bool
	IsPrefix(s *discordgo.Session, m *discordgo.MessageCreate) bool
	GetCommandData() *discordgo.ApplicationCommand
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
	Prefix(s *discordgo.Session, m *discordgo.MessageCreate)
}

type command struct {
	cmd *discordgo.ApplicationCommand
}

func (c command) GetName() string {
	if c.cmd != nil {
		return c.cmd.Name
	}
	return "UnknownCommand"
}

func (c command) IsCommand(i *discordgo.InteractionCreate) bool {
	if c.cmd != nil {
		return i.ApplicationCommandData().Name == c.cmd.Name
	}
	return false
}

func (c command) IsPrefix(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	return m.Author.ID != s.State.User.ID && strings.HasPrefix(m.Content, fmt.Sprintf("%s%s", config.BotConfig().Prefix(), c.GetName()))
}

func (c command) GetCommandData() *discordgo.ApplicationCommand {
	return c.cmd
}

func (c command) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Printf("unimplement command handler")
}

func (c command) Prefix(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Printf("unimplement prefix handler")
}

func interactionRespond(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		fmt.Printf("Error while interaction respond %v\n", err)
	}
}

func interactionRespondEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		fmt.Printf("Error while interaction respond %v\n", err)
	}
}
