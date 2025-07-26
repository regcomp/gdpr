package caching

import (
	"sync"
	"time"
)

type LocalCache struct {
	cookieHashes *cookieHashStore
	nonces       *nonceStore
	sessions     *sessionStore
	requests     *requestStore
}

func createLocalCache() *LocalCache {
	return &LocalCache{
		cookieHashes: createCookieHashStore(),
		nonces:       createNonceStore(),
		sessions:     createSessionStore(),
		requests:     createRequestStore(),
	}
}

type cookieHashStore struct {
	mu  sync.RWMutex
	dat []byte
}

func createCookieHashStore() *cookieHashStore {
	return &cookieHashStore{
		mu:  sync.RWMutex{},
		dat: make([]byte, 0, 128),
	}
}

type nonceStore struct {
	mu  sync.RWMutex
	dat map[string]string
}

func createNonceStore() *nonceStore {
	return &nonceStore{
		mu:  sync.RWMutex{},
		dat: make(map[string]string),
	}
}

type sessionStore struct {
	mu  sync.RWMutex
	dat map[string][]byte
}

func createSessionStore() *sessionStore {
	return &sessionStore{
		mu:  sync.RWMutex{},
		dat: make(map[string][]byte),
	}
}

type requestStore struct {
	mu  sync.RWMutex
	dat map[string][]byte
}

func createRequestStore() *requestStore {
	return &requestStore{
		mu:  sync.RWMutex{},
		dat: make(map[string][]byte),
	}
}

func (lc *LocalCache) NonceAdd(nonce string) error {
	lc.nonces.mu.Lock()
	defer lc.nonces.mu.Unlock()

	timestamp := time.UTC.String()
	lc.nonces.dat[nonce] = timestamp

	return nil
}

func (lc *LocalCache) SessionAdd(sessionID string, data []byte) error {
	lc.sessions.mu.Lock()
	defer lc.sessions.mu.Unlock()

	lc.sessions.dat[sessionID] = data
	return nil
}

func (lc *LocalCache) SessionGet(sessionID string) ([]byte, error) {
	lc.sessions.mu.RLock()
	defer lc.sessions.mu.RUnlock()

	sessionData, _ := lc.sessions.dat[sessionID]

	return sessionData, nil
}

func (lc *LocalCache) CookieHashesGet() ([]byte, error) {
	lc.cookieHashes.mu.RLock()
	defer lc.cookieHashes.mu.RUnlock()

	return lc.cookieHashes.dat, nil
}

func (lc *LocalCache) CookieHashesSet(data []byte) error {
	lc.cookieHashes.mu.Lock()
	defer lc.cookieHashes.mu.Unlock()

	lc.cookieHashes.dat = data
	return nil
}

func (lc *LocalCache) RequestAdd(requestID string, data []byte) error {
	lc.requests.mu.Lock()
	defer lc.requests.mu.Unlock()

	lc.requests.dat[requestID] = data

	return nil
}

func (lc *LocalCache) RequestRetrieve(requestID string) ([]byte, error) {
	lc.requests.mu.RLock()
	defer lc.requests.mu.RUnlock()

	data, _ := lc.requests.dat[requestID]

	return data, nil
}
