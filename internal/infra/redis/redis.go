package redis

import (
	"neko-bot/internal/config"

	"github.com/redis/go-redis/v9"
)

const (
	WatchedThreadIds string = "watched_thread_ids"
	WatchedZennUsers string = "watched_zenn_users"
	NoticeChannel    string = "notice_channel"
	IpaSecurityAlert string = "ipa_security_alert"
	IpaNoticeChannel string = "ipa_notice_channel"
	LastAction       string = "last_action"
)

var instance *Cache

type Cache struct {
	client *redis.Client
}

func Client() *redis.Client {
	if instance == nil {
		conf := config.RedisConfig()
		client := redis.NewClient(&redis.Options{
			Addr:     conf.RedisHost() + ":" + conf.RedisPort(),
			Password: conf.RedisPassword(),
			DB:       0,
		})

		instance = &Cache{
			client: client,
		}
	}

	return instance.client
}
