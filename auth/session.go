package auth

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/regcomp/gdpr/caching"
)

type SessionStore struct {
	cache caching.IServiceCache
}

type SessionData struct {
	//
}

func CreateSessionStore(serviceCache caching.IServiceCache) *SessionStore {
	return &SessionStore{cache: serviceCache}
}

func (ss *SessionStore) CreateSession() (string, error) {
	id := generateSessionID()
	sessionData := &SessionData{}
	sessionBytes, err := json.Marshal(sessionData)
	if err != nil {
		return "", err
	}

	ss.cache.AddSession(id, sessionBytes)
	return id, nil
}

func (ss *SessionStore) UpdateSession(data *SessionData) error {
	//
	return nil
}

func (ss *SessionStore) GetSession(sessionID string) (*SessionData, error) {
	sessionBytes, err := ss.cache.GetSession(sessionID)
	if err != nil {
		return nil, err
	}

	sessionData := &SessionData{}
	err = json.Unmarshal(sessionBytes, sessionData)
	if err != nil {
		return nil, err
	}

	return sessionData, nil
}

func (ss *SessionStore) DeleteSession(sessionID string) {
	//
}

func generateSessionID() string {
	return uuid.NewString()
}
