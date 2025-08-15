package caching

import (
	"encoding/hex"
	"net/http"
	"time"

	sc "github.com/gorilla/securecookie"
	"github.com/regcomp/gdpr/internal/config"
)

type NonceManager struct {
	cache INonceStash
}

func CreateNonceManager(cache INonceStash) *NonceManager {
	return &NonceManager{
		cache: cache,
	}
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
