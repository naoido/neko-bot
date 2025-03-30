package bot

import "neko-bot/discord/neko"

var (
	bot    *neko.Bot
	config *neko.Config
)

func Start() error {
	var err error

	config = neko.GetConfig()
	bot, err = neko.New(*config)
	return err
}

func Update() error {
	var err error
	config, err = neko.ReloadConfig()
	if err != nil {
		return err
	}

	bot.UpdateBot(*config, true)

	return nil
}

func Stop() error {
	return bot.Stop()
}
