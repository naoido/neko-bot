package command

import (
	"github.com/bwmarrin/discordgo"
	"neko-bot/internal/errors"
)

func interactionRespond(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		errors.Catch(err, err.Error())
	}
}
