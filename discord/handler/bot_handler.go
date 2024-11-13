package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"neko-bot/discord/command"
)

var (
	commands           []*command.Detail
	registeredCommands []*discordgo.ApplicationCommand
)

func Add(command *command.Detail) {
	commands = append(commands, command)
}

func handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	for _, model := range commands {
		model.Handler(s, i)
	}
}

func prefixHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, model := range commands {
		model.Prefix(s, m)
	}
}

func RegisterCommands(session *discordgo.Session) error {
	// Add parent handler
	session.AddHandler(handler)
	session.AddHandler(prefixHandler)

	// Add slash command data
	for _, model := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, model.GetCommandData().GuildID, model.GetCommandData())
		if err != nil {
			return err
		}
		registeredCommands = append(registeredCommands, cmd)
		fmt.Printf("\rRegisterd command %s\r\n", model.GetCommandData().Name)
	}
	return nil
}

func RemoveCommands(session *discordgo.Session) error {
	for _, model := range registeredCommands {
		err := session.ApplicationCommandDelete(session.State.User.ID, model.GuildID, model.ID)
		if err != nil {
			return err
		}
		fmt.Printf("\rRemove command %s\r\n", model.Name)
	}
	return nil
}

func GetRegisteredCommands() []*discordgo.ApplicationCommand {
	return registeredCommands
}
