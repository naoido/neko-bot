package handler

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"neko-bot/internal/infra/redis"
)

func onReaction(_ *discordgo.Session, c *discordgo.MessageReactionAdd) {
	userId := c.UserID

	redis.UpdateLastActionTime(context.Background(), userId)
}
