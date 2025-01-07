package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"neko-bot/discord/handler"
	"neko-bot/internal/errors"
)

func RegisterEvent() error {
	bot.Session().AddHandler(forumNotify)
	return handler.RegisterCommands(bot.Session())
}

func forumNotify(s *discordgo.Session, c *discordgo.MessageCreate) {
	Channel, err := s.Channel(c.ChannelID)
	errors.Catch(err, "cannot get channel")

	if Channel.Type == discordgo.ChannelTypeGuildPublicThread && Channel.ParentID == "ForumID" {
		name := c.Author.Username
		url := fmt.Sprintf("https://discord.com/channels/%s/%s", c.ChannelID, c.Message.ID)
		emoji, err := s.State.Emoji(c.GuildID, "EmojiID")
		errors.Catch(err, "cannot get emoji")

		messages, err := s.ChannelMessages(c.ChannelID, 100, "", "", "")
		errors.Catch(err, "cannot get messages")

		Forum, err := s.Channel(Channel.ParentID)
		errors.Catch(err, "cannot get forum channel")

		tagList := Forum.AvailableTags

		var tags string
		for _, tagID := range Channel.AppliedTags {
			for _, tag := range tagList {
				if tag.ID == tagID {
					tags += fmt.Sprintf("#%s ", tag.Name)
				}
			}
		}

		if len(messages) <= 1 {
			_, err = s.ChannelMessageSend("ChannelID", name+"さんが記事が投稿したよ! すごい！！！！"+emoji.MessageFormat()+"\n"+url+tags)
			errors.Catch(err, "cannot send message")
		}
	}
}
