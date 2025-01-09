package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type Say struct {
	Detail
}

func NewSay(name string, prefix *string) *Say {
	say := &Say{
		Detail: Detail{
			name:   name,
			prefix: prefix,
		},
	}

	say.Detail.Command = say

	return say
}

func (say *Say) GetName() string {
	return say.Detail.name
}

func (say *Say) GetCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        say.GetName(),
		Description: "say command",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "message to say",
				Required:    true,
			},
		},
	}
}

func (say *Say) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !say.Detail.isCommand(i) {
		return
	}

	option := i.ApplicationCommandData().Options[0].StringValue()
	if option == "" {
		interactionRespond(s, i, "空文字は送信できません")
	}
	interactionRespond(s, i, option)
}

func (say *Say) Prefix(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !say.Detail.isPrefix(s, m) {
		return
	}

	if strings.HasPrefix(m.Content, fmt.Sprintf("%s%s ", *say.prefix, say.name)) {
		s.ChannelMessageSend(m.ChannelID, strings.Replace(m.Content, fmt.Sprintf("%s%s ", *say.prefix, say.name), "", 1))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s%s { 内容 }", *say.prefix, say.name))
	}
}
