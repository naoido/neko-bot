package command

import (
	"fmt"
	"neko-bot/internal/errors"
	"neko-bot/redis"
	"time"

	"github.com/bwmarrin/discordgo"
)

type LastAction struct {
	Detail
}

func NewLastAction(name string, prefix *string) *LastAction {
	lastAction := &LastAction{
		Detail: Detail{
			name:   name,
			prefix: prefix,
		},
	}

	lastAction.Detail.Command = lastAction

	return lastAction
}

func (l *LastAction) GetName() string {
	return l.Detail.name
}

func (l *LastAction) GetCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        l.GetName(),
		Description: "get last user action time.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user_id",
				Description: "user id",
			},
		},
	}
}

func (l *LastAction) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !l.Detail.isCommand(i) {
		return
	}

	options := i.ApplicationCommandData().Options
	userId := i.Member.User.ID
	if len(options) != 0 {
		userId = options[0].StringValue()
	}

	lastActionTime, _ := redis.GetLastActionTime(userId)
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		errors.Catch(err, "Faild to load time location")
		interactionRespond(s, i, "タイムゾーン取得時にエラーが発生しました。")
		return
	}

	jstTime := time.Unix(lastActionTime, 0).In(jst)
	formattedTime := jstTime.Format("2006-01-02 15:04:05")
	interactionRespond(s, i, fmt.Sprintf("<@!%s>さんの最後の行動は`%s`です。", userId, formattedTime))
}

func (l *LastAction) Prefix(s *discordgo.Session, m *discordgo.MessageCreate) {
	return
}
