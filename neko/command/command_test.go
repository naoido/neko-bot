package command

import (
    "testing"
    

    "github.com/bwmarrin/discordgo"
)

type MockSession struct {
    ChannelID string
    Message   string
}

func (m *MockSession) ChannelMessageSend(channelID string, message string) (*discordgo.Message, error) {
    m.ChannelID = channelID
    m.Message = message
    return &discordgo.Message{}, nil
}



func TestUmbrellaCommand(t *testing.T) {
    mockSession := &MockSession{}
    messageCreate := &discordgo.MessageCreate{
        Message: &discordgo.Message{
            Content:   "/umbrella",
            ChannelID: "testChannel",
        },
    }

    UmbrellaCommand(mockSession, messageCreate)

    if mockSession.ChannelID != "testChannel" {
        t.Errorf("expected channel ID to be 'testChannel', got %s", mockSession.ChannelID)
    }

    expectedMessage := "MrsGreenApple「umbrella」: https://www.youtube.com/watch?v=RnBQela7oyE"
    if mockSession.Message != expectedMessage {
        t.Errorf("expected message to be '%s', got %s", expectedMessage, mockSession.Message)
    }
}
