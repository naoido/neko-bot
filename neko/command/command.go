package command

import (
    "fmt"

    "github.com/bwmarrin/discordgo"
)

// Define an interface that includes the ChannelMessageSend method
type SessionInterface interface {
    ChannelMessageSend(channelID string, message string) (*discordgo.Message, error)
}



// UmbrellaCommand handles the /umbrella command
func UmbrellaCommand(s SessionInterface, m *discordgo.MessageCreate) {
    if m.Content == "/umbrella" {
        url := "https://www.youtube.com/watch?v=RnBQela7oyE"
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("MrsGreenApple「umbrella」: %s", url))
    }
}
