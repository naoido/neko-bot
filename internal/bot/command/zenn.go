package command

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"neko-bot/internal/client/zenn"
	"neko-bot/internal/infra/redis"
	"strings"
	"time"
)

type Zenn struct {
	command
}

var notifier *zenn.ArticleNotifier

func NewZenn(s *discordgo.Session) *Zenn {
	notifier = zenn.NewArticleNotifier()
	notifier.Init()
	go notifier.Start()

	go func() {
		for {
			message := []string{"🆕新しい記事を発見しました！"}
			for _, newArticle := range <-notifier.NewArticleChan {
				message = append(message, fmt.Sprintf("%s ) `%s`", newArticle.User.Name, newArticle.Title), fmt.Sprintf("https://zenn.dev/%s/articles/%s", newArticle.User.Username, newArticle.Slug))
			}
			noticeChannel, err := redis.Client().Get(context.Background(), redis.NoticeChannel).Result()
			if err != nil || noticeChannel == "" {
				fmt.Println("redis error:", err)
				return
			}
			s.ChannelMessageSend(noticeChannel, strings.Join(message, "\n"))
		}
	}()

	return &Zenn{
		command{
			&discordgo.ApplicationCommand{
				Name:        "zenn",
				Description: "zenn command",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "add",
						Description: "add watching user",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "user_id",
								Description: "user id",
								Required:    true,
							},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "remove",
						Description: "remove watching user",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "user_id",
								Description: "user id",
								Required:    true,
							},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "list",
						Description: "list watching user",
					},
				},
			},
		},
	}
}

func (z *Zenn) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch options[0].Name {
	case "add":
		newUser := options[0].Options[0].StringValue()
		err := notifier.AddUser(ctx, newUser)
		if err != nil {
			interactionRespond(s, i, fmt.Sprintf("エラーが発生しました: %v", err))
		} else {
			interactionRespond(s, i, fmt.Sprintf("新しくウォッチリストに追加しました！ https://zenn.dev/%s", newUser))
		}
	case "remove":
		newUser := options[0].Options[0].StringValue()
		found, err := notifier.RemoveUser(ctx, newUser)
		if err != nil {
			interactionRespond(s, i, fmt.Sprintf("エラーが発生しました: %v", err))
		} else if found {
			interactionRespond(s, i, "ユーザーが見つかりませんでした")
		} else {
			interactionRespond(s, i, fmt.Sprintf("削除しました: %s", newUser))
		}
	case "list":
		users := notifier.WatchUsers
		if len(users) == 0 {
			interactionRespond(s, i, "まだ何も登録されていません。")
		}
		content := "登録済みリスト📝\n"
		for i, user := range users {
			content += fmt.Sprintf("%d) %s\n", i+1, user)
		}

		interactionRespond(s, i, content)
	}
}
