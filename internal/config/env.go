package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"neko-bot/internal/util/zr"
	"os"
	"strings"
)

type config struct {
	token         string
	prefix        string
	stage         string
	activeMessage string
	status        string
	developers    []string
	chatGPTKey    string
	redisHost     string
	redisPort     string
	redisPassword string
}

var conf *config

func init() {
	fmt.Println("start loading environment variables")
	stage := zr.OrDef(os.Getenv("STAGE"), "dev")
	if stage != "prod" {
		err := godotenv.Overload(fmt.Sprintf("env/%s.env", stage))
		// env/{stage}.envが読み込めない場合はpanic
		if err != nil {
			fmt.Println("could not load .env file")
			panic(err)
		}
	}

	conf = &config{
		token:         "Bot " + os.Getenv("DISCORD_TOKEN"),
		prefix:        zr.OrDef(os.Getenv("DISCORD_PREFIX"), "!"),
		stage:         zr.OrDef(os.Getenv("DISCORD_STAGE"), "prod"),
		activeMessage: zr.OrDef(os.Getenv("DISCORD_ACTIVITY_MESSAGE"), "Just chilling..."),
		status:        zr.OrDef(os.Getenv("DISCORD_STATUS_TYPE"), "online"),
		developers:    strings.Split(zr.OrDef(os.Getenv("DEVELOPERS"), ""), ","),
		chatGPTKey:    zr.OrDef(os.Getenv("CHATGPT_API_KEY"), "Bot"),
		redisHost:     zr.OrDef(os.Getenv("REDIS_HOST"), "localhost"),
		redisPort:     zr.OrDef(os.Getenv("REDIS_PORT"), "6379"),
		redisPassword: zr.OrDef(os.Getenv("REDIS_PASSWORD"), ""),
	}

	fmt.Println("end loading environment variables")
}
