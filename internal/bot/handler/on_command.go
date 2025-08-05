package handler

import (
	"github.com/bwmarrin/discordgo"
)

func onCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	for _, cmd := range commands {
		if cmd.IsCommand(i) {
			cmd.Handler(s, i)
			return
		}
	}
}
