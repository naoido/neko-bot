package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"neko-bot/internal/client/chatgpt"
)

func isMentionToMe(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	for _, user := range m.Mentions {
		if user != nil && user.ID == s.State.User.ID {
			return true
		}
	}
	return false
}

func onMention(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !isMentionToMe(s, m) || s.State.User.ID == m.Author.ID {
		return
	}

	err := s.ChannelTyping(m.ChannelID)
	if err != nil {
		fmt.Println("cant start typing", err)
		return
	}
	res, err := chatgpt.GetResponse(m.ContentWithMentionsReplaced())
	if err != nil {
		fmt.Println("ChatGPT GetResponse error", err)
		_, err = s.ChannelMessageSend(m.ChannelID, "エラーが発生しました。")
		return
	}
	for _, choice := range res.Choices {
		_, err = s.ChannelMessageSend(m.ChannelID, choice.Message.Content)
		if err != nil {
			fmt.Println("cant sent message to channel", err)
			return
		}
	}
}
