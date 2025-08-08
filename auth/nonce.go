package auth

import (
	"encoding/hex"
	"net/http"
	"time"

	sc "github.com/gorilla/securecookie"
	"github.com/regcomp/gdpr/caching"
	"github.com/regcomp/gdpr/config"
)

const (
	nonceTimeout         = 15
	nonceCleanupInterval = 5
)

type NonceManager struct {
	cache caching.IServiceCache
}

func CreateNonceStash(serviceCache caching.IServiceCache) *NonceManager {
	store := &NonceManager{cache: serviceCache}

	return store
}

func (ns *NonceManager) Generate() string {
	bytes := sc.GenerateRandomKey(32)

	nonce := hex.EncodeToString(bytes)
	timestamp := time.UTC.String()

	ns.cache.StashNonce(nonce, timestamp)

	return nonce
}

func (ns *NonceManager) Validate(r *http.Request) bool {
	nonce := r.FormValue(config.FormValueNonce)
	if nonce == "" {
		nonce = r.Header.Get(config.HeaderNonceToken) // ajax
	}

	// TODO: ACTUALLY VALIDATE

	return true
}
