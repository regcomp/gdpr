package caching

import (
	"encoding/json"
	"net/http"

	sc "github.com/gorilla/securecookie"
	"github.com/regcomp/gdpr/internal/auth"
	"github.com/regcomp/gdpr/internal/config"
)

type CookieManager struct {
	cache ICookieHashesStore
	keys  *sc.SecureCookie
}

func CreateCookieManager(cache ICookieHashesStore) *CookieManager {
	return &CookieManager{
		cache: cache,
		keys:  nil,
	}
}

func (cm *CookieManager) CreateAccessCookie(accessToken string) (*http.Cookie, error) {
	keys, err := cm.getKeys()
	if err != nil {
		return nil, err
	}
	return auth.CreateCookie(
		config.CookieNameAccessToken,
		accessToken,
		keys,
		// TODO: Configure
	)
}

func (cm *CookieManager) GetAccessToken(r *http.Request) (string, error) {
	keys, err := cm.getKeys()
	if err != nil {
		return "", err
	}
	return auth.GetTokenFromCookie(r, config.CookieNameAccessToken, keys)
}

func (cm *CookieManager) CreateRefreshCookie(refreshToken string) (*http.Cookie, error) {
	keys, err := cm.getKeys()
	if err != nil {
		return nil, err
	}
	return auth.CreateCookie(
		config.CookieNameRefreshToken,
		refreshToken,
		keys,
		// TODO: Configure
		func(c *http.Cookie) {
			c.Path = config.RouterAuthPathPrefix + config.EndpointAuthRenewToken
		},
	)
}

func (cm *CookieManager) GetRefreshToken(r *http.Request) (string, error) {
	keys, err := cm.getKeys()
	if err != nil {
		return "", err
	}
	return auth.GetTokenFromCookie(r, config.CookieNameRefreshToken, keys)
}

func (cm *CookieManager) CreateSessionCookie(sessionID string) (*http.Cookie, error) {
	keys, err := cm.getKeys()
	if err != nil {
		return nil, err
	}
	return auth.CreateCookie(
		config.CookieNameSessionId,
		sessionID,
		keys,
		// TODO: Configure
	)
}

func (cm *CookieManager) GetSessionID(r *http.Request) (string, error) {
	keys, err := cm.getKeys()
	if err != nil {
		return "", err
	}
	return auth.GetTokenFromCookie(r, config.CookieNameSessionId, keys)
}

func (cm *CookieManager) DestroyAllCookies(w http.ResponseWriter, r *http.Request) {
	for _, name := range config.CookieNames {
		cookie, err := r.Cookie(name)
		if err != nil {
			continue
		}
		auth.DestroyCookie(w, cookie)
	}
}

func (cm *CookieManager) getKeys() (*sc.SecureCookie, error) {
	if cm.keys == nil {
		freshHashes := &cookieHashes{}

		cookieHashesBytes, err := cm.cache.GetCookieHashes()
		if err != nil {
			return nil, err
		}

		if len(cookieHashesBytes) == 0 {
			freshHashes = newCookieHashes()
			cookieHashesBytes, err = json.Marshal(freshHashes)
			if err != nil {
				return nil, err
			}
			err = cm.cache.SetCookieHashes(cookieHashesBytes)
			if err != nil {
				return nil, err
			}
		} else {
			err = json.Unmarshal(cookieHashesBytes, freshHashes)
			if err != nil {
				return nil, err
			}
		}

		cm.keys = sc.New(freshHashes.Hash, freshHashes.Block)
	}

	return cm.keys, nil
}

type cookieHashes struct {
	Hash  []byte `json:"hash"`
	Block []byte `json:"block"`
}

func newCookieHashes() *cookieHashes {
	return &cookieHashes{
		Hash:  sc.GenerateRandomKey(64),
		Block: sc.GenerateRandomKey(32),
	}
}
