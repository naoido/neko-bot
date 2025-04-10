package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"neko-bot/discord/neko"
)

const (
	WatchedThreadIds string = "watched_thread_ids"
	NoticeChannel    string = "notice_channel"
	IpaSecurityAlert string = "ipa_security_alert"
	IpaNoticeChannel string = "ipa_notice_channel"
)

var instance *Cache

type Cache struct {
	ctx    context.Context
	client *redis.Client
}

func Client() *redis.Client {
	if instance == nil {
		config := neko.GetConfig()
		client := redis.NewClient(&redis.Options{
			Addr:     config.RedisHost + ":" + config.RedisPort,
			Password: "",
			DB:       0,
		})

		instance = &Cache{
			ctx:    context.Background(),
			client: client,
		}
	}

	return instance.client
}

func Context() context.Context {
	return instance.ctx
}
