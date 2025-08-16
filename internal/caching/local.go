package caching

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type LocalCache struct {
	mu  sync.RWMutex
	dat map[string][]byte
	// constants generated at creation to avoid collisions
	cookieHashesKey string
}

func CreateLocalCache() *LocalCache {
	return &LocalCache{
		mu:              sync.RWMutex{},
		dat:             make(map[string][]byte),
		cookieHashesKey: uuid.New().String(),
	}
}

func (lc *LocalCache) StashRequest(id string, data []byte) error {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.dat[id] = data

	return nil
}

func (lc *LocalCache) RetrieveRequest(id string) ([]byte, error) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	data, ok := lc.dat[id]
	if !ok {
		return nil, fmt.Errorf("could not find cached request")
	}

	delete(lc.dat, id)

	return data, nil
}

func (lc *LocalCache) StashNonce(nonce, timestamp string) error {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.dat[nonce] = []byte(timestamp)

	return nil
}

func (lc *LocalCache) RetrieveNonce(nonce string) (string, error) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	timestamp, ok := lc.dat[nonce]
	if !ok {
		return "", fmt.Errorf("could not find cached request")
	}

	delete(lc.dat, nonce)

	return string(timestamp), nil
}

func (lc *LocalCache) AddSession(id string, data []byte) error {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.dat[id] = data
	return nil
}

func (lc *LocalCache) GetSession(id string) ([]byte, error) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	sessionData, ok := lc.dat[id]
	if !ok {
		return nil, fmt.Errorf("could not find session=%s", id)
	}

	return sessionData, nil
}

func (lc *LocalCache) SetCookieHashes(data []byte) error {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.dat[lc.cookieHashesKey] = data
	return nil
}

func (lc *LocalCache) GetCookieHashes() ([]byte, error) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	return bytes.Clone(lc.dat[lc.cookieHashesKey]), nil
}
