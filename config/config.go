package config

import (
	"github.com/regcomp/gdpr/constants"
)

var attrs = []string{
	constants.ConfigServiceURLKey,
	constants.ConfigDefaultPortKey,
	constants.ConfigSessionDurationKey,
}

type IConfigStore interface {
	GetServiceURL() string
	GetDefaultPort() string
	GetSessionDuration() string
}

type LocalConfigStore struct {
	// mu    sync.RWMutex // May need in the future
	attrs map[string]string
}

func NewLocalConfigStore() *LocalConfigStore {
	return &LocalConfigStore{
		attrs: make(map[string]string),
	}
}

func (cs *LocalConfigStore) InitializeStore(getters ...func(string) string) {
	for _, getter := range getters {
		for _, attr := range attrs {
			cs.attrs[attr] = getter(attr)
		}
	}
}

func (cs *LocalConfigStore) GetServiceURL() string {
	return cs.attrs[constants.ConfigServiceURLKey]
}

func (cs *LocalConfigStore) GetDefaultPort() string {
	return cs.attrs[constants.ConfigDefaultPortKey]
}

func (cs *LocalConfigStore) GetSessionDuration() string {
	return cs.attrs[constants.ConfigSessionDurationKey]
}
