package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"neko-bot/internal/errors"
	"os"
)

func RegisterEvent() error {
	bot.Session().AddHandler(forumNotify)
	return nil
}

func forumNotify(s *discordgo.Session, c *discordgo.MessageCreate) {
	channel, err := s.Channel(c.ChannelID)
	errors.Catch(err, "cannot get channel")

	if channel.Type == discordgo.ChannelTypeGuildPublicThread && channel.ParentID == os.Getenv("FORUM_CHANNEL_ID") {
		name := c.Author.Username
		url := fmt.Sprintf("https://discord.com/channels/%s/%s", c.ChannelID, c.Message.ID)
		emoji, err := s.State.Emoji(c.GuildID, os.Getenv("EMOJI_TENSAI_ID"))
		errors.Catch(err, "cannot get emoji")

		messages, err := s.ChannelMessages(c.ChannelID, 10, "", "", "")
		errors.Catch(err, "cannot get messages")

		forum, err := s.Channel(channel.ParentID)
		errors.Catch(err, "cannot get forum channel")

		tagList := forum.AvailableTags

		var tags string
		for _, tagID := range channel.AppliedTags {
			for _, tag := range tagList {
				if tag.ID == tagID {
					tags += fmt.Sprintf("#%s ", tag.Name)
				}
			}
		}

		var imageURL string
		if len(c.Message.Attachments) > 0 {
			imageURL = c.Message.Attachments[0].URL
		} else {
			imageURL = "https://random-image-pepebigotes.vercel.app/api/random-image"
		}
		if len(messages) <= 1 {
			_, err = s.ChannelMessageSend(os.Getenv("SEND_CHANNEL_ID"), name+"さんが記事を投稿したよ! すごい！！！！"+emoji.MessageFormat())
			errors.Catch(err, "cannot send message")
			_, err = s.ChannelMessageSendEmbed(os.Getenv("SEND_CHANNEL_ID"), &discordgo.MessageEmbed{
				Title:       url,
				Description: tags,
				Color:       0xf54900,
				Author: &discordgo.MessageEmbedAuthor{
					Name: name,
				},
				Image: &discordgo.MessageEmbedImage{
					URL: imageURL,
				},
			})
			errors.Catch(err, "cannot send message")
		}
	}
}
