package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"neko-bot/internal/errors"
	"time"
)

type Voice struct {
	Detail
}

func NewVoice(name string, prefix *string) *Voice {
	voice := &Voice{
		Detail: Detail{
			name:   name,
			prefix: prefix,
		},
	}

	voice.Detail.Command = voice

	return voice
}

func (v *Voice) GetName() string {
	return v.Detail.name
}

func (v *Voice) GetCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        v.GetName(),
		Description: "voice command",
	}
}

func (v *Voice) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !v.Detail.isCommand(i) {
		return
	}

	v.handleVoiceCommand(s, i.GuildID, i.Member.User.ID, i.ChannelID)
}

func (v *Voice) Prefix(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !v.Detail.isPrefix(s, m) {
		return
	}

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
		errors.Catch(err, "cannot send message")
		return
	}

	_, err := s.ChannelVoiceJoin(guildId, userstate.ChannelID, false, false)
	errors.Catch(err, "cannot join voice channel")

	_, err = s.ChannelMessageSend(chId, "やあ!😄")
	errors.Catch(err, "cannot send voice command")

	time.Sleep(5 * time.Second)

	err = s.VoiceConnections[guildId].Disconnect()
	errors.Catch(err, "cannot disconnect voice channel")
}
