package auth

import (
	"context"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	sc "github.com/gorilla/securecookie"
)

const (
	nonceTimeout         = 15
	nonceCleanupInterval = 5
)

type INonceStore interface {
	Generate() string
	Validate(*http.Request) bool
}

type NonceStore struct {
	mu     sync.RWMutex
	nonces map[string]time.Time
	ctx    context.Context
	cancel context.CancelFunc
}

func CreateNonceStore() *NonceStore {
	ctx, cancel := context.WithCancel(context.Background())
	store := &NonceStore{
		nonces: make(map[string]time.Time),
		ctx:    ctx,
		cancel: cancel,
	}

	go store.cleanupDaemon()
	return store
}

func (ns *NonceStore) Generate() string {
	bytes := sc.GenerateRandomKey(32)

	nonce := hex.EncodeToString(bytes)

	ns.mu.Lock()
	ns.nonces[nonce] = time.Now().Add(nonceTimeout * time.Minute)
	ns.mu.Unlock()

	return nonce
}

func (ns *NonceStore) Validate(r *http.Request) bool {
	nonce := r.FormValue("nonce")
	if nonce == "" {
		nonce = r.Header.Get("CSRF-Nonce") // ajax
	}
	ns.mu.Lock()
	defer ns.mu.Unlock()

	expiry, exists := ns.nonces[nonce]
	if !exists {
		return false
	}

	if time.Now().After(expiry) {
		delete(ns.nonces, nonce)
		return false
	}

	delete(ns.nonces, nonce) // number used once!
	return true
}

func (ns *NonceStore) cleanupDaemon() {
	ticker := time.NewTicker(nonceCleanupInterval * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ns.ctx.Done():
			return
		case <-ticker.C:
			ns.cleanup()
		}
	}
}

func (ns *NonceStore) cleanup() {
	ns.mu.RLock()
	now := time.Now()
	toDelete := make([]string, 0)
	for nonce, expiry := range ns.nonces {
		if now.After(expiry) {
			toDelete = append(toDelete, nonce)
		}
	}
	ns.mu.RUnlock()

	if len(toDelete) > 0 {
		ns.mu.Lock()
		for _, nonce := range toDelete {
			delete(ns.nonces, nonce)
		}
		ns.mu.Unlock()
	}
}
