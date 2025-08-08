package caching

import (
	"fmt"

	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/secrets"
)

/*

Things that should live in the cache:
  - auth provider credentials
  - all database credentials
    - support for multiple databases
  - cookie hashes
    - hash key
    - block key
  - session store
  - nonce store
  - request store
  - logger/audit credentials
    - where should logs/audits be sent
  -

*/

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
		return createLocalCache(), nil
	default:
		return nil, fmt.Errorf("unknown cache type=%s", cfg.CacheType)
	}
}
