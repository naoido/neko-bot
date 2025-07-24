package redis

import (
	"context"
	"fmt"
	"time"
)

func GetLastActionTime(ctx context.Context, userId string) (int64, error) {
	key := fmt.Sprintf("%s:%s", LastAction, userId)

	value, err := Client().Get(ctx, key).Int64()
	if value == 0 {
		UpdateLastActionTime(ctx, userId)
		value, err = Client().Get(ctx, key).Int64()
	}

	return value, err
}

func UpdateLastActionTime(ctx context.Context, userId string) {
	key := fmt.Sprintf("%s:%s", LastAction, userId)
	now := time.Now().Unix()

	if err := Client().Set(ctx, key, now, 0).Err(); err != nil {
		fmt.Println("redis set error", err)
	}
}
