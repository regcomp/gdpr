package sessions

import (
	"crypto/rand"
	"log"

	gsessions "github.com/gorilla/sessions"
)

type SessionManager struct {
	AccessStore  gsessions.CookieStore
	RefreshStore gsessions.CookieStore
}

func CreateSessionManager() *SessionManager {
	return &SessionManager{
		AccessStore:  createAccessStore(generateSecureKey(), generateSecureKey()),
		RefreshStore: createRefreshStore(generateSecureKey(), generateSecureKey()),
	}
}

func createAccessStore(authenticationKey, encryptionKey []byte) gsessions.CookieStore {
	// TODO: configure settings
	return *gsessions.NewCookieStore(authenticationKey, encryptionKey)
}

func createRefreshStore(authenticationKey, encryptionKey []byte) gsessions.CookieStore {
	// TODO: contigure settings. This should only be sent on a specific route
	return *gsessions.NewCookieStore(authenticationKey, encryptionKey)
}

func generateSecureKey() []byte {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		log.Fatal("Failed to generate session key:", err)
	}
	return key
}
