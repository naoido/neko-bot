package neko

import (
	"neko-bot/internal/errors"
	"neko-bot/internal/zr"
	"os"
	"strings"
)

type Config struct {
	token         string
	Prefix        string
	stage         string
	activeMessage string
	status        string
	developers    []string
	ChatgptKey    string
	RedisHost     string
	RedisPort     string
}

var config *Config

func init() {
	c, err := loadConfig(getStage())
	errors.Catch(err, "failed to load env file")

	config = c
}

func getStage() string {
	return zr.OrDef(os.Getenv("STAGE"), "dev")
}

func GetConfig() *Config {
	return config
}

func ReloadConfig() (*Config, error) {
	return loadConfig(getStage())
}

func loadConfig(stage string) (*Config, error) {
	//if stage != "prod" {
	//	err := godotenv.Overload(fmt.Sprintf("env/%s.env", stage))
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	token := "Bot " + os.Getenv("DISCORD_TOKEN")
	prefix := zr.OrDef(os.Getenv("DISCORD_COMMAND_PREFIX"), "!")
	statusType := zr.OrDef(os.Getenv("DISCORD_STATUS_TYPE"), "online")
	activeMessage := zr.OrDef(os.Getenv("DISCORD_ACTIVITY_MESSAGE"), "Just chilling...")
	developers := strings.Split(zr.OrDef(os.Getenv("DEVELOPERS"), ""), ",")
	chatgptKey := zr.OrDef(os.Getenv("CHATGPT_API_KEY"), "Bot")
	redisHost := zr.OrDef(os.Getenv("REDIS_HOST"), "localhost")
	redisPort := zr.OrDef(os.Getenv("REDIS_PORT"), "6379")

	return &Config{
		token:         token,
		Prefix:        prefix,
		stage:         stage,
		activeMessage: activeMessage,
		status:        statusType,
		developers:    developers,
		ChatgptKey:    chatgptKey,
		RedisHost:     redisHost,
		RedisPort:     redisPort,
	}, nil
}
