package handler

import "github.com/bwmarrin/discordgo"

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, hdl := range commands {
		if hdl.IsPrefix(s, m) {
			hdl.Prefix(s, m)
			break
		}
	}

	if isMentionToMe(s, m) && s.State.User.ID != m.Author.ID {
		onMention(s, m)
		return
	}
}
