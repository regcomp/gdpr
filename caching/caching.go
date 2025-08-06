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

// NOTE: This will likely be a large and cluttered interface
type IServiceCache interface {
	// Nonce Handling
	NonceAdd(string) error

	// Session Handling
	SessionAdd(string, []byte) error
	SessionGet(string) ([]byte, error)

	// Cookie Hashes
	CookieHashesGet() ([]byte, error)
	CookieHashesSet([]byte) error

	// Requests
	RequestAdd(string, []byte) error
	RequestRetrieve(string) ([]byte, error)
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
