package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type Ping struct {
	command
}

func NewPing() *Ping {
	return &Ping{
		command{
			&discordgo.ApplicationCommand{
				Name:        "ping",
				Description: "ping command",
			},
		},
	}
}

func (p *Ping) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	interactionRespond(s, i, "pong")
}

func (p *Ping) Prefix(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "pong")
	if err != nil {
		fmt.Println("Error sending pong: ", err)
	}
}
