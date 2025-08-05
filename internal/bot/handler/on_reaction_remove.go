package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"slices"
)

func onReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	handleUnPin(s, r)
}

func handleUnPin(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if !slices.Contains([]string{"📌", "📍"}, r.Emoji.Name) {
		return
	}

	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		fmt.Println("Error getting message", err)
		return
	}

	if !message.Pinned {
		return
	}

	channel, err := s.Channel(message.ChannelID)
	if err != nil {
		fmt.Println("Error getting channel", err)
		return
	}

	if !channel.IsThread() {
		return
	}

	permittedReactionUsers := make([]*discordgo.User, 0)
	for _, targetReaction := range []string{"📌", "📍"} {
		reactedUsers, err := s.MessageReactions(r.ChannelID, r.MessageID, targetReaction, 100, "", "")
		if err != nil {
			fmt.Println("Error getting reacted users", err)
			return
		}

		for _, user := range reactedUsers {
			if channel.OwnerID == user.ID {
				permittedReactionUsers = append(permittedReactionUsers, user)
			}
		}
	}

	// 権限を持ってる人が全員ピンを外したら
	if len(permittedReactionUsers) == 0 {
		err = s.ChannelMessageUnpin(message.ChannelID, message.ID)
		if err != nil {
			fmt.Println("Error unpinning message", err)
		}
	}
}
