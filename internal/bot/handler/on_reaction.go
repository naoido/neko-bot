package handler

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"neko-bot/internal/infra/redis"
	"neko-bot/internal/util/discordutil"
	"slices"
)

func onReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	userId := r.UserID

	pinHandler(s, r)
	redis.UpdateLastActionTime(context.Background(), userId)
}

func pinHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// :pushpin:📌 or :round_pushpin:📍で反応する
	if !slices.Contains([]string{"📌", "📍"}, r.Emoji.Name) {
		return
	}

	channel, err := s.State.Channel(r.ChannelID)
	if err != nil {
		fmt.Println("Error getting channel", err)
		return
	}

	if !channel.IsThread() {
		return
	}

	member, err := s.GuildMember(channel.GuildID, r.UserID)
	if err != nil {
		fmt.Println("Error getting guild member", err)
		return
	}

	if channel.OwnerID == r.UserID || discordutil.HasAdminPermission(member) {
		message, err := s.ChannelMessage(channel.ID, r.MessageID)
		if err != nil {
			fmt.Println("Error getting message", err)
			return
		}

		err = s.ChannelMessagePin(message.ChannelID, message.ID)
		if err != nil {
			fmt.Println("Error pining message", err)
		}
	}
}
