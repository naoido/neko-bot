package command

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

type Say struct {
	command
}

func NewSay() *Say {
	return &Say{
		command{
			&discordgo.ApplicationCommand{
				Name:        "say",
				Description: "say command",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "message",
						Description: "message to say",
						Required:    true,
					},
				},
			},
		},
	}
}

func (say *Say) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	option := i.ApplicationCommandData().Options[0].StringValue()
	if option == "" {
		interactionRespond(s, i, "空文字は送信できません")
	}
	interactionRespondEphemeral(s, i, "送信しました！")
	s.ChannelMessageSend(i.ChannelID, option)
}

func (say *Say) Prefix(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")
	content := strings.Join(args[1:], " ")
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, content)
	} else {
		s.ChannelMessageSend(m.ChannelID, content)
	}
}
