package config

type Bot interface {
	Token() string
	Prefix() string
	ActiveMessage() string
	Status() string
}

func BotConfig() Bot {
	return conf
}

func (conf *config) Token() string {
	return conf.token
}

func (conf *config) Prefix() string {
	return conf.prefix
}

func (conf *config) ActiveMessage() string {
	return conf.activeMessage
}

func (conf *config) Status() string {
	return conf.status
}
