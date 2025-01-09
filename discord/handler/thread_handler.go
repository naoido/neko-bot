package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"neko-bot/internal/errors"
	"neko-bot/redis"
	"slices"
)

func threadHandler(session *discordgo.Session, c *discordgo.ThreadCreate) {
	threadId := c.Channel.ID
	threadChannel, err := session.Channel(threadId)
	if err != nil || threadChannel == nil {
		return
	}

	watched := redis.Client().SMembers(redis.Context(), redis.WatchedThreadIds).Val()

	if threadChannel.Type == discordgo.ChannelTypeGuildPublicThread && slices.Contains(watched, threadChannel.ParentID) {
		noticeChannel, err := redis.Client().Get(redis.Context(), redis.NoticeChannel).Result()
		if err != nil {
			errors.Catch(err, err.Error())
			return
		}
		messages, err := session.ChannelMessages(threadId, 1, "", "", "")
		errors.Catch(err, "cannot get messages")

		if len(messages) == 0 {
			fmt.Println("No messages in the thread")
			return
		}

		username := messages[0].Author.Username
		messageURL := fmt.Sprintf("https://discord.com/channels/%s/%s", threadId, messages[0].ID)

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

		_, err = session.ChannelMessageSend(noticeChannel, username+"さんが記事を投稿したよ! すごい！！！！<a:tensai:1325578494683512902>")
		errors.Catch(err, "cannot send message")

		_, err = session.ChannelMessageSendEmbed(noticeChannel, &discordgo.MessageEmbed{
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
