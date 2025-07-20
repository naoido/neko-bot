package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"neko-bot/internal/errors"
	"neko-bot/internal/zenn"
	"neko-bot/redis"
	"strings"
)

type Zenn struct {
	Detail
}

var notifier *zenn.ArticleNotifier

func NewZenn(name string, prefix *string, s *discordgo.Session) *Zenn {
	z := &Zenn{
		Detail: Detail{
			name:   name,
			prefix: prefix,
		},
	}

	notifier = zenn.NewArticleNotifier()
	notifier.Init()
	go notifier.Start()

	go func() {
		for {
			message := []string{"🆕新しい記事を発見しました！"}
			for _, newArticle := range <-notifier.NewArticleChan {
				message = append(message, fmt.Sprintf("%s ) `%s`", newArticle.User.Name, newArticle.Title), fmt.Sprintf("https://zenn.dev/%s/articles/%s", newArticle.User.Username, newArticle.Slug))
			}
			noticeChannel, err := redis.Client().Get(redis.Context(), redis.NoticeChannel).Result()
			if err != nil {
				errors.Catch(err, err.Error())
				return
			}
			s.ChannelMessageSend(noticeChannel, strings.Join(message, "\n"))
		}
	}()

	z.Detail.Command = z

	return z
}

func (z *Zenn) GetName() string {
	return z.Detail.name
}

func (z *Zenn) GetCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        z.GetName(),
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
	}
}

func (z *Zenn) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !z.Detail.isCommand(i) {
		return
	}
	options := i.ApplicationCommandData().Options

	switch options[0].Name {
	case "add":
		newUser := options[0].Options[0].StringValue()
		err := notifier.AddUser(newUser)
		if err != nil {
			interactionRespond(s, i, fmt.Sprintf("エラーが発生しました: %v", err))
		} else {
			interactionRespond(s, i, fmt.Sprintf("新しくウォッチリストに追加しました！ https://zenn.dev/%s", newUser))
		}
	case "remove":
		newUser := options[0].Options[0].StringValue()
		found, err := notifier.RemoveUser(newUser)
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

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "pong",
		},
	})
	errors.Catch(err, "cannot respond to ping")
}

func (z *Zenn) Prefix(s *discordgo.Session, m *discordgo.MessageCreate) {}
