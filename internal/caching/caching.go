package caching

import (
	"fmt"

	"github.com/regcomp/gdpr/internal/config"
	"github.com/regcomp/gdpr/internal/secrets"
)

type IServiceCache interface {
	INonceStash
	IRequestStash
	ISessionStore
	ICookieHashesStore
}

type IRequestStash interface {
	StashRequest(string, []byte) error
	RetrieveRequest(string) ([]byte, error)
}

type INonceStash interface {
	StashNonce(string, string) error
	RetrieveNonce(string) (string, error)
}

type ISessionStore interface {
	AddSession(string, []byte) error
	GetSession(string) ([]byte, error)
}

type ICookieHashesStore interface {
	SetCookieHashes([]byte) error
	GetCookieHashes() ([]byte, error)
}

func CreateServiceCache(
	cfg *config.ServiceCacheConfig,
	secrets *secrets.ServiceCacheSecrets,
) (IServiceCache, error) {
	switch cfg.CacheType {
	case config.ValueLocalType:
		return CreateLocalCache(), nil
	default:
		return nil, fmt.Errorf("unknown cache type=%s", cfg.CacheType)
	}
}
