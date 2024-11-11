package manager

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"neko-bot/internal/errors"
	"neko-bot/neko"
)

var (
	Commands        = make([]*discordgo.ApplicationCommand, 0)
	CommandHandlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
)

var (
	s = neko.GetDiscord()
)

var (
	commandList            = Commands
	registeredCommandsList = make([]*discordgo.ApplicationCommand, 0)
)

func addCommand(command *discordgo.ApplicationCommand, fn func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	_, exists := CommandHandlers[command.Name]
	err := fmt.Errorf("[%s] This command name was duplcated", command.Name)
	if exists {
		errors.CatchAndPanic(err, "Command already exists")
	}

	CommandHandlers[command.Name] = fn
	Commands = append(Commands, command)
}

func registerCommand() error {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commandList))
	for i, v := range commandList {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			return err
		}
		registeredCommands[i] = cmd
	}
	return nil
}

func removeCommands() error {
	for _, v := range registeredCommandsList {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
