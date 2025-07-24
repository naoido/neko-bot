package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

type Voice struct {
	command
}

func NewVoice() *Voice {
	return &Voice{
		command{
			&discordgo.ApplicationCommand{
				Name:        "voice",
				Description: "voice command",
			},
		},
	}
}

func (v *Voice) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	v.handleVoiceCommand(s, i.GuildID, i.Member.User.ID, i.ChannelID)
}

func (v *Voice) Prefix(s *discordgo.Session, m *discordgo.MessageCreate) {
	v.handleVoiceCommand(s, m.GuildID, m.Author.ID, m.ChannelID)
}

func (v *Voice) handleVoiceCommand(s *discordgo.Session, guildId, authorId, chId string) {
	userstate, _ := s.State.VoiceState(guildId, authorId)
	if userstate == nil {
		fmt.Println("User is not in a voice channel")
		return
	}

	if vc, ok := s.VoiceConnections[guildId]; ok && vc != nil {
		_, err := s.ChannelMessageSend(chId, "ちゃんと部屋みてねw")
		if err != nil {
			fmt.Println("Error sending pong: ", err)
		}
		return
	}

	_, err := s.ChannelVoiceJoin(guildId, userstate.ChannelID, false, false)
	if err != nil {
		fmt.Println("Error joining voice channel: ", err)
	}

	_, err = s.ChannelMessageSend(chId, "やあ!😄")
	if err != nil {
		fmt.Println("Error sending pong: ", err)
	}

	time.Sleep(5 * time.Second)

	err = s.VoiceConnections[guildId].Disconnect()
	if err != nil {
		fmt.Println("Error disconnecting voice channel: ", err)
	}
}
