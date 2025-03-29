package handler

import "github.com/bwmarrin/discordgo"

func RegisterHandlers(session *discordgo.Session) {
	session.AddHandler(threadHandler)
}
