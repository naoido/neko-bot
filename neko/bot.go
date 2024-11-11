package neko

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"neko-bot/internal/errors"
	"neko-bot/internal/zr"
	"os"
)

var (
	discord *discordgo.Session
	config  map[string]string
)

func Start(stage string) {
	var err error
	config, err = getConfig(stage)
	errors.CatchAndPanic(err, "failed to get config")

	discord, err = discordgo.New(config["token"])
	errors.CatchAndPanic(err, "failed to create session of discord bot")

	err = discord.Open()
	errors.CatchAndPanic(err, "failed to open discord bot")

	UpdateBot(false)

	const nekoBot = `
          _____                    _____                    _____                   _______                   _____                   _______               _____
         /\    \                  /\    \                  /\    \                 /::\    \                 /\    \                 /::\    \             /\    \
        /::\____\                /::\    \                /::\____\               /::::\    \               /::\    \               /::::\    \           /::\    \
       /::::|   |               /::::\    \              /:::/    /              /::::::\    \             /::::\    \             /::::::\    \          \:::\    \
      /:::::|   |              /::::::\    \            /:::/    /              /::::::::\    \           /::::::\    \           /::::::::\    \          \:::\    \
     /::::::|   |             /:::/\:::\    \          /:::/    /              /:::/~~\:::\    \         /:::/\:::\    \         /:::/~~\:::\    \          \:::\    \
    /:::/|::|   |            /:::/__\:::\    \        /:::/____/              /:::/    \:::\    \       /:::/__\:::\    \       /:::/    \:::\    \          \:::\    \
   /:::/ |::|   |           /::::\   \:::\    \      /::::\    \             /:::/    / \:::\    \     /::::\   \:::\    \     /:::/    / \:::\    \         /::::\    \
  /:::/  |::|   | _____    /::::::\   \:::\    \    /::::::\____\________   /:::/____/   \:::\____\   /::::::\   \:::\    \   /:::/____/   \:::\____\       /::::::\    \
 /:::/   |::|   |/\    \  /:::/\:::\   \:::\    \  /:::/\:::::::::::\    \ |:::|    |     |:::|    | /:::/\:::\   \:::\ ___\ |:::|    |     |:::|    |     /:::/\:::\    \
/:: /    |::|   /::\____\/:::/__\:::\   \:::\____\/:::/  |:::::::::::\____\|:::|____|     |:::|    |/:::/__\:::\   \:::|    ||:::|____|     |:::|    |    /:::/  \:::\____\
\::/    /|::|  /:::/    /\:::\   \:::\   \::/    /\::/   |::|~~~|~~~~~      \:::\    \   /:::/    / \:::\   \:::\  /:::|____| \:::\    \   /:::/    /    /:::/    \::/    /
 \/____/ |::| /:::/    /  \:::\   \:::\   \/____/  \/____|::|   |            \:::\    \ /:::/    /   \:::\   \:::\/:::/    /   \:::\    \ /:::/    /    /:::/    / \/____/
         |::|/:::/    /    \:::\   \:::\    \            |::|   |             \:::\    /:::/    /     \:::\   \::::::/    /     \:::\    /:::/    /    /:::/    /
         |::::::/    /      \:::\   \:::\____\           |::|   |              \:::\__/:::/    /       \:::\   \::::/    /       \:::\__/:::/    /    /:::/    /
         |:::::/    /        \:::\   \::/    /           |::|   |               \::::::::/    /         \:::\  /:::/    /         \::::::::/    /     \::/    /
         |::::/    /          \:::\   \/____/            |::|   |                \::::::/    /           \:::\/:::/    /           \::::::/    /       \/____/
         /:::/    /            \:::\    \                |::|   |                 \::::/    /             \::::::/    /             \::::/    /
        /:::/    /              \:::\____\               \::|   |                  \::/____/               \::::/    /               \::/____/
        \::/    /                \::/    /                \:|   |                   ~~                      \::/____/                 ~~
         \/____/                  \/____/                  \|___|                                            ~~

`
	fmt.Print(nekoBot)
}

func Stop() {
	fmt.Println("\r\u001B[00;31mShutting down...\u001B[00m")
	errors.Catch(discord.Close(), "could not close discord bot")
}

func UpdateBot(reload bool) {
	if reload {
		var err error
		stage := os.Getenv("STAGE")
		config, err = getConfig(stage)
		errors.Catch(err, "could not get config")

		err = discord.Close()
		errors.Catch(err, "could not close discord")

		err = discord.Open()
		errors.CatchAndPanic(err, "could not open discord")
	}
	err := discord.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: config["status_type"],
		Activities: []*discordgo.Activity{
			{
				Name: config["activity_message"],
				Type: discordgo.ActivityTypeGame,
			},
		},
	})
	errors.Catch(err, "\rfailed to update status of discord bot\r")
}

func getConfig(stage string) (map[string]string, error) {
	config := make(map[string]string)
	err := godotenv.Overload(fmt.Sprintf("env/%s.env", stage))
	if err != nil {
		return nil, err
	}

	err = os.Setenv("STAGE", stage)
	errors.Catch(err, "\rfailed to set STAGE")

	config["token"] = "Bot " + os.Getenv("DISCORD_TOKEN")
	config["activity_message"] = zr.OrDef(os.Getenv("DISCORD_ACTIVITY_MESSAGE"), "Just chilling...")
	config["status_type"] = zr.OrDef(os.Getenv("DISCORD_STATUS_TYPE"), "online")

	return config, nil
}

func GetDiscord() *discordgo.Session {
	return discord
}
