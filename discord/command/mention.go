package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"neko-bot/chatgpt"
	"neko-bot/internal/errors"
)

type Mention struct {
	Detail
}

func NewMention() *Mention {
	mention := &Mention{
		Detail: Detail{
			name:   "",
			prefix: nil,
		},
	}
	mention.Detail.Command = mention

	return mention
}

func (mention *Mention) GetName() string                               { return mention.name }
func (mention *Mention) GetCommandData() *discordgo.ApplicationCommand { return nil }

func isMention(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	for _, user := range m.Mentions {
		if user != nil && user.ID == s.State.User.ID {
			return true
		}
	}
	return false
}

func idToMention(id string) string { return fmt.Sprintf("<@!%s>", id) }

func (mention *Mention) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {}

func (mention *Mention) Prefix(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !isMention(s, m) || s.State.User.ID == m.Author.ID {
		return
	}

	s.ChannelTyping(m.ChannelID)
	res, err := chatgpt.GetResponse(m.ContentWithMentionsReplaced())
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, "エラーが発生しました。")
		return
	}
	for _, choice := range res.Choices {
		_, err = s.ChannelMessageSend(m.ChannelID, choice.Message.Content)
	}
	errors.Catch(err, "Failed to send message")
}
