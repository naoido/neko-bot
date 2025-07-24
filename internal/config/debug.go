package config

type Debug interface {
	Stage() string
}

func DebugConfig() Debug {
	return conf
}

func (conf *config) Stage() string {
	return conf.stage
}
