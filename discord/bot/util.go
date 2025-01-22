package bot

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

func SendMessage(channelID string, message string) error {
	if channelID == "" {
		return errors.New("channelID is empty")
	}

	s := bot.Session()
	if s == nil {
		return errors.New("session is nil")
	}

	_, err := s.ChannelMessageSend(channelID, message)
	return err
}

func SendMessageEmbed(channelID string, embed *discordgo.MessageEmbed) error {
	if channelID == "" {
		return errors.New("channelID is empty")
	}

	s := bot.Session()
	if s == nil {
		return errors.New("session is nil")
	}

	_, err := s.ChannelMessageSendEmbed(channelID, embed)
	return err
}
