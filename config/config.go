package config

/*
The idea for this package is that the .env file should not be
edited that much because it has sensitive contents. This file contains
general server configuration.
*/

const (
	hostPath        = "HOST_PATH"
	defaultPort     = "DEFAULT_PORT"
	sessionDuration = "SESSION_DURATION"
)

var attrs = []string{
	hostPath,
	defaultPort,
	sessionDuration,
}

type IConfigStore interface {
	GetHostPath() string
	GetSessionDuration() string
}

type ConfigStore struct {
	attrs map[string]string
}

func NewConfigStore() *ConfigStore {
	return &ConfigStore{}
}

func (cs *ConfigStore) InitializeStore(getters ...func(string) string) {
	for _, getter := range getters {
		for _, attr := range attrs {
			cs.attrs[attr] = getter(attr)
		}
	}
}

func (cs *ConfigStore) GetHostPath() string {
	return cs.attrs[hostPath]
}

func (cs *ConfigStore) GetSessionDuration() string {
	return cs.attrs[sessionDuration]
}

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
