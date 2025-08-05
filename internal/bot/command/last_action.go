package command

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"neko-bot/internal/infra/redis"
	"time"
)

type LastAction struct {
	command
}

func NewLastAction() *LastAction {
	return &LastAction{
		command{
			&discordgo.ApplicationCommand{
				Name:        "action",
				Description: "get last user action time.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "user_id",
						Description: "user id",
					},
				},
			},
		},
	}
}

func (l *LastAction) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	userId := i.Member.User.ID
	if len(options) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userId = options[0].StringValue()
	lastActionTime, _ := redis.GetLastActionTime(ctx, userId)
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Printf("load jst err: %v\n", err)
		interactionRespond(s, i, "タイムゾーン取得時にエラーが発生しました。")
		return
	}

	jstTime := time.Unix(lastActionTime, 0).In(jst)
	formattedTime := jstTime.Format("2006-01-02 15:04:05")
	interactionRespond(s, i, fmt.Sprintf("<@!%s>さんの最後の行動は`%s`です。", userId, formattedTime))
}
