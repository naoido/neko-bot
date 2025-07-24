package config

type ChatGPT interface {
	ChatGPTKey() string
}

func ChatGPTConfig() ChatGPT {
	return conf
}

func (conf *config) ChatGPTKey() string {
	return conf.chatGPTKey
}
