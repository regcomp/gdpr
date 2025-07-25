package caching

import (
	"sync"

	"github.com/regcomp/gdpr/secrets"
)

type LocalCache struct {
	secretsStore secrets.ISecretStore

	cookieHashes *cookieHashStore
	nonces       *nonceStore
	sessions     *sessionStore
	requests     *requestStore
}

func createLocalCache(secretsStore secrets.ISecretStore) *LocalCache {
	return &LocalCache{
		secretsStore: secretsStore,

		cookieHashes: createCookieHashStore(),
		nonces:       createNonceStore(),
		sessions:     createSessionStore(),
		requests:     createRequestStore(),
	}
}

type cookieHashStore struct {
	mu           sync.RWMutex
	cookieHashes []byte
}

func createCookieHashStore() *cookieHashStore {
	return &cookieHashStore{
		mu:           sync.RWMutex{},
		cookieHashes: make([]byte, 0, 128),
	}
}

type nonceStore struct {
	mu     sync.RWMutex
	nonces map[string]string
}

func createNonceStore() *nonceStore {
	return &nonceStore{
		mu:     sync.RWMutex{},
		nonces: make(map[string]string),
	}
}

type sessionStore struct {
	mu       sync.RWMutex
	sessions map[string][]byte
}

func createSessionStore() *sessionStore {
	return &sessionStore{
		mu:       sync.RWMutex{},
		sessions: make(map[string][]byte),
	}
}

type requestStore struct {
	mu       sync.RWMutex
	sessions map[string][]byte
}

func createRequestStore() *requestStore {
	return &requestStore{
		mu:       sync.RWMutex{},
		sessions: make(map[string][]byte),
	}
}
func (lc *LocalCache) NonceAdd(nonce string) {}

func (lc *LocalCache) SessionAdd(sessionID string, data []byte)    {}
func (lc *LocalCache) SessionGet(sessionID string) ([]byte, error) { return nil, nil }

func (lc *LocalCache) CookieHashesGet() []byte                { return nil }
func (lc *LocalCache) CookieHashesSet([]byte)                 {}
func (lc *LocalCache) RequestAdd(string, []byte) error        { return nil }
func (lc *LocalCache) RequestRetrieve(string) ([]byte, error) { return nil, nil }
func (lc *LocalCache) DatabaseGetConfig() any
