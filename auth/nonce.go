package auth

import (
	"encoding/hex"
	"net/http"

	sc "github.com/gorilla/securecookie"
	"github.com/regcomp/gdpr/caching"
)

const (
	nonceTimeout         = 15
	nonceCleanupInterval = 5
)

type NonceStore struct {
	cache caching.IServiceCache
}

func CreateNonceStore(serviceCache caching.IServiceCache) *NonceStore {
	store := &NonceStore{
		cache: serviceCache,
	}

	return store
}

func (ns *NonceStore) Generate() string {
	bytes := sc.GenerateRandomKey(32)

	nonce := hex.EncodeToString(bytes)

	ns.cache.NonceAdd(nonce)

	return nonce
}

func (ns *NonceStore) Validate(r *http.Request) bool {
	nonce := r.FormValue("nonce")
	if nonce == "" {
		nonce = r.Header.Get("CSRF-Nonce") // ajax
	}

	return true
}
