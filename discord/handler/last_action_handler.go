package handler

import (
	"neko-bot/redis"

	"github.com/bwmarrin/discordgo"
)

func messageHandler(session *discordgo.Session, c *discordgo.MessageCreate) {
	userId := c.Author.ID

	redis.UpdateLastActionTime(userId)
}

func reactionHandler(session *discordgo.Session, c *discordgo.MessageReactionAdd) {
	userId := c.UserID

	redis.UpdateLastActionTime(userId)
}
