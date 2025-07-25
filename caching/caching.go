package caching

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

// NOTE: This will likely be a large and cluttered interface
type IServiceCache interface {
	// Nonce Handling
	NonceAdd(string)

	// Session Handling
	SessionAdd(string, []byte)
	SessionGet(string) ([]byte, error)

	// Cookie Hashes
	CookieHashesGet() []byte
	CookieHashesSet([]byte)

	// Requests
	RequestAdd(string, []byte) error
	RequestRetrieve(string) ([]byte, error)

	// Database
	DatabaseGetConfig() any
}

func CreateServiceCache(secretStore secrets.ISecretStore, cacheType string) (IServiceCache, error) {
	switch cacheType {
	default:
		return createLocalCache(secretStore), nil
	}
}
