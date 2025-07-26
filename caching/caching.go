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

const localCacheType = "LOCAL"

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
	config *config.ServiceCacheConfig,
	secrets *secrets.ServiceCacheSecrets,
) (IServiceCache, error) {
	switch config.CacheType {
	case localCacheType:
		return createLocalCache(), nil
	default:
		return nil, fmt.Errorf("unknown cache type=%s", config.CacheType)
	}
}
