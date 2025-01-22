package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"neko-bot/redis"
)

type Setting struct {
	Detail
}

func NewSetting(name string, prefix *string) *Setting {
	setting := &Setting{
		Detail: Detail{
			name:   name,
			prefix: prefix,
		},
	}

	setting.Detail.Command = setting

	return setting
}

func (setting *Setting) GetName() string {
	return setting.Detail.name
}

func (setting *Setting) GetCommandData() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        setting.GetName(),
		Description: "guild setting command",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Name:        "thread",
				Description: "thread",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "add",
						Description: "add to watch a thread",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "thread_id",
								Description: "thread id",
								Required:    true,
							},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "remove",
						Description: "remove from watch a thread",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "thread_id",
								Description: "thread id",
								Required:    true,
							},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "list",
						Description: "list threads",
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "notice",
						Description: "set notice for create thread",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "channel_id",
								Description: "channel id",
								Required:    false,
							},
						},
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "ipa",
				Description: "ipa",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "channel_id",
						Description: "channel id",
						Required:    false,
					},
				},
			},
		},
	}
}

func (setting *Setting) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !setting.isCommand(i) {
		return
	}
	options := i.ApplicationCommandData().Options

	switch options[0].Name {
	case "thread":
		options = options[0].Options
		switch options[0].Name {
		case "add":
			threadId := options[0].Options[0].StringValue()
			err := redis.Client().SAdd(redis.Context(), redis.WatchedThreadIds, threadId).Err()
			if err != nil {
				interactionRespond(s, i, fmt.Sprintf("エラーが発生しました: %v", err))
			} else {
				interactionRespond(s, i, fmt.Sprintf("新しくウォッチリストに追加しました！ <#%v>", threadId))
			}
		case "remove":
			threadId := options[0].Options[0].StringValue()
			err := redis.Client().SRem(redis.Context(), redis.WatchedThreadIds, threadId).Err()
			if err != nil {
				interactionRespond(s, i, fmt.Sprintf("エラーが発生しました: %v", err))
			} else {
				interactionRespond(s, i, fmt.Sprintf("<#%v>を削除しました", threadId))
			}
		case "list":
			threads := redis.Client().SMembers(redis.Context(), redis.WatchedThreadIds).Val()
			if len(threads) == 0 {
				interactionRespond(s, i, "まだ何も登録されていません。")
			}
			content := "登録済みリスト📝\n"
			for j, thread := range threads {
				content += fmt.Sprintf("%d) <#%v>(%v)\n", j+1, thread, thread)
			}

			interactionRespond(s, i, content)
		case "notice":
			if options[0].Options == nil || len(options[0].Options) == 0 {
				noticeChannel := redis.Client().Get(redis.Context(), redis.NoticeChannel).Val()
				interactionRespond(s, i, fmt.Sprintf("現在は <#%v> に設定されています。", noticeChannel))
			} else {
				newNoticeChannel := options[0].Options[0].StringValue()
				if newNoticeChannel == "" {
					interactionRespond(s, i, "設定するにはChannelIDを指定してください")
					return
				}
				err := redis.Client().Set(redis.Context(), redis.NoticeChannel, newNoticeChannel, 0).Err()
				if err != nil {
					interactionRespond(s, i, "エラーが発生しました")
					return
				}
				interactionRespond(s, i, fmt.Sprintf("<#%v>を新しい通知チャンネルに設定しました", newNoticeChannel))
			}
		}
	case "ipa":
		if options[0].Options == nil || len(options[0].Options) == 0 {
			noticeChannel := redis.Client().Get(redis.Context(), redis.IpaNoticeChannel).Val()
			interactionRespond(s, i, fmt.Sprintf("現在は <#%v> に設定されています。", noticeChannel))
		} else {
			newNoticeChannel := options[0].Options[0].StringValue()
			if newNoticeChannel == "" {
				interactionRespond(s, i, "設定するにはChannelIDを指定してください")
				return
			}
			err := redis.Client().Set(redis.Context(), redis.IpaNoticeChannel, newNoticeChannel, 0).Err()
			if err != nil {
				interactionRespond(s, i, "エラーが発生しました")
				return
			}
			interactionRespond(s, i, fmt.Sprintf("<#%v>を新しい通知チャンネルに設定しました", newNoticeChannel))
		}
	}
}

func (setting *Setting) Prefix(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !setting.isPrefix(s, m) {
		return
	}
	s.ChannelMessageSendReply(m.ChannelID, "このコマンドはPrefixでは使用できません", m.Reference())
}
