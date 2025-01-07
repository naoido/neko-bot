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

func forumNotify(session *discordgo.Session, threadCreate *discordgo.ThreadCreate) {
	threadID := threadCreate.Channel.ID
	threadChannel, err := session.Channel(threadID)
	errors.Catch(err, "cannot get channel")

	if threadChannel.Type == discordgo.ChannelTypeGuildPublicThread && threadChannel.ParentID == os.Getenv("FORUM_CHANNEL_ID") {
		messages, err := session.ChannelMessages(threadID, 1, "", "", "")
		errors.Catch(err, "cannot get messages")

		if len(messages) == 0 {
			fmt.Println("No messages in the thread")
			return
		}

		username := messages[0].Author.Username
		messageURL := fmt.Sprintf("https://discord.com/channels/%s/%s", threadID, messages[0].ID)
		emoji, err := session.State.Emoji(threadCreate.Channel.GuildID, os.Getenv("EMOJI_TENSAI_ID"))
		errors.Catch(err, "cannot get emoji")

		forumChannel, err := session.Channel(threadChannel.ParentID)
		errors.Catch(err, "cannot get forum channel")

		tagList := forumChannel.AvailableTags

		var tags string
		for _, tagID := range threadChannel.AppliedTags {
			for _, tag := range tagList {
				if tag.ID == tagID {
					tags += fmt.Sprintf("#%s ", tag.Name)
				}
			}
		}

		var imageURL string
		if len(messages[0].Attachments) > 0 {
			imageURL = messages[0].Attachments[0].URL
		} else {
			imageURL = "https://random-image-pepebigotes.vercel.app/api/random-image"
		}

		_, err = session.ChannelMessageSend(os.Getenv("SEND_CHANNEL_ID"), username+"さんが記事を投稿したよ! すごい！！！！"+emoji.MessageFormat())
		errors.Catch(err, "cannot send message")

		_, err = session.ChannelMessageSendEmbed(os.Getenv("SEND_CHANNEL_ID"), &discordgo.MessageEmbed{
			Title:       messageURL,
			Description: tags,
			Color:       0xf54900,
			Author: &discordgo.MessageEmbedAuthor{
				Name: username,
			},
			Image: &discordgo.MessageEmbedImage{
				URL: imageURL,
			},
		})
		errors.Catch(err, "cannot send message")
	}
}
