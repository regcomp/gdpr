package auth

import "github.com/google/uuid"

type SessionStore struct {
	SessionIDToHMACKey map[string]string
}

func CreateSessionStore() *SessionStore {
	return &SessionStore{
		SessionIDToHMACKey: make(map[string]string),
	}
}

func generateSessionID() string {
	return uuid.NewString()
}
