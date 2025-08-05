package config

type Redis interface {
	RedisHost() string
	RedisPort() string
	RedisPassword() string
}

func RedisConfig() Redis {
	return conf
}

func (conf *config) RedisHost() string {
	return conf.redisHost
}

func (conf *config) RedisPort() string {
	return conf.redisPort
}

func (conf *config) RedisPassword() string {
	return conf.redisPassword
}
