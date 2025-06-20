package config

/*
The idea for this package is that the .env file should not be
edited that much because it has sensitive contents. This file contains
general server configuration.
*/

type Config struct {
	Port string // 16 bit uint
	//
}

func LoadConfig(getenv func(string) string) *Config {
	config := &Config{}

	config.Port = getenv("CONFIG_PORT")

	return config
}
