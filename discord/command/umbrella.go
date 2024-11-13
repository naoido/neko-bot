package command

import (
    "github.com/bwmarrin/discordgo"
    "neko-bot/internal/errors"
)

type Umbrella struct {
    Model
}

func NewUmbrella(name string, prefix *string) *Umbrella {
    umbrella := &Umbrella{
        Model: Model{
            Detail: Detail{
                name:   name,
                prefix: prefix,
            },
        },
    }

    umbrella.Model.Command = umbrella

    return umbrella
}

func (u *Umbrella) GetName() string {
    return u.Model.Detail.name
}

func (u *Umbrella) GetCommandData() *discordgo.ApplicationCommand {
    return &discordgo.ApplicationCommand{
        Name:        u.GetName(),
        Description: "Returns the YouTube URL for MRS's Umbrella",
    }
}

func (u *Umbrella) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
    if (!u.Model.Detail.isCommand(i)) {
        return
    }

    err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: "https://www.youtube.com/watch?v=6Ul-wFnquno", // MRSのUmbrellaのURLに置き換えてください
        },
    })
    errors.Catch(err, "cannot respond to umbrella command")
}

func (u *Umbrella) Prefix(s *discordgo.Session, m *discordgo.MessageCreate) {
    if (!u.Model.Detail.isPrefix(s, m)) {
        return
    }

    _, err := s.ChannelMessageSend(m.ChannelID, "https://www.youtube.com/watch?v=6Ul-wFnquno") // MRSのUmbrellaのURLに置き換えてください
    errors.Catch(err, "cannot send umbrella message")
}