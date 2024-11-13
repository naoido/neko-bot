package command

import (
	"github.com/bwmarrin/discordgo"
	"neko-bot/internal/errors"
)

type Ping struct {
	Detail
}

func NewPing(name string, prefix *string) *Ping {
	ping := &Ping{
		Detail: Detail{
			name:   name,
			prefix: prefix,
		},
	}

	ping.Detail.Command = ping

	return ping
}

func (p *Ping) GetName() string {
	return p.Detail.name
}

func (p *Ping) GetCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        p.GetName(),
		Description: "ping command",
	}
}

func (p *Ping) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !p.Detail.isCommand(i) {
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "pong",
		},
	})
	errors.Catch(err, "cannot respond to ping")
}

func (p *Ping) Prefix(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !p.Detail.isPrefix(s, m) {
		return
	}

	_, err := s.ChannelMessageSend(m.ChannelID, "pong")
	errors.Catch(err, "cannot send ping")
}
