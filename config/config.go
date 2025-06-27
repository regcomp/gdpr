package config

/*
The idea for this package is that the .env file should not be
edited that much because it has sensitive contents. This file contains
general server configuration.
*/

type Config struct {
	Port    string // 16 bit uint
	TLSKey  string
	TLSCert string
}

func LoadConfig(getenv func(string) string) *Config {
	config := &Config{}

	config.Port = getenv("CONFIG_PORT")

	config.TLSKey = "/auth/local_certs/server.key"
	config.TLSCert = "/auth/local_certs/server.crt"

	return config
}
