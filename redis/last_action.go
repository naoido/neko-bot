package redis

import (
	"fmt"
	"neko-bot/internal/errors"
	"time"
)

func GetLastActionTime(userId string) (int64, error) {
	key := fmt.Sprintf("%s:%s", LastAction, userId)

	value, err := Client().Get(Context(), key).Int64()
	if value == 0 {
		UpdateLastActionTime(userId)
		value, err = Client().Get(Context(), key).Int64()
	}

	return value, err
}

func UpdateLastActionTime(userId string) {
	key := fmt.Sprintf("%s:%s", LastAction, userId)
	now := time.Now().Unix()

	if err := Client().Set(Context(), key, now, 0).Err(); err != nil {
		errors.Catch(err, "Faild to update last user action time.")
	}
}
