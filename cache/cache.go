package cache

import (
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
	ServiceCache()
}

func CreateServiceCache(secretStore secrets.ISecretStore, cacheType string) (IServiceCache, error) {
	switch cacheType {
	default:
		return createInMemoryCache(), nil
	}
}

type InMemoryCache struct {
	//
}

func createInMemoryCache() *InMemoryCache {
	return &InMemoryCache{}
}

func (imc *InMemoryCache) ServiceCache() {}
