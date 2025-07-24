package handler

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"neko-bot/internal/infra/redis"
	"slices"
	"time"
)

func onThreadCreate(session *discordgo.Session, c *discordgo.ThreadCreate) {
	threadId := c.Channel.ID
	threadChannel, err := session.Channel(threadId)
	if err != nil || threadChannel == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	watched := redis.Client().SMembers(ctx, redis.WatchedThreadIds).Val()

	if threadChannel.Type == discordgo.ChannelTypeGuildPublicThread && slices.Contains(watched, threadChannel.ParentID) {
		noticeChannel, err := redis.Client().Get(ctx, redis.NoticeChannel).Result()
		if err != nil {
			fmt.Println(err)
			return
		}
		messages, err := session.ChannelMessages(threadId, 1, "", "", "")
		if err != nil {
			fmt.Println(err)
		}

		// フォーラムでない場合は無視
		if len(messages) == 0 {
			fmt.Println("this is not a foram")
			return
		}

		username := messages[0].Author.Username
		messageURL := fmt.Sprintf("https://discord.com/channels/%s/%s", threadId, messages[0].ID)

		forumChannel, err := session.Channel(threadChannel.ParentID)
		if err != nil {
			fmt.Println(err)
			return
		}

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
		if err != nil {
			fmt.Println(err)
			return
		}

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
		if err != nil {
			fmt.Println(err)
		}
	}
}
