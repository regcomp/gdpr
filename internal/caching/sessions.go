package caching

import (
	"encoding/json"

	"github.com/google/uuid"
)

type SessionManager struct {
	cache ISessionStore
}

type SessionData struct {
	//
}

func CreateSessionManager(cache ISessionStore) *SessionManager {
	return &SessionManager{cache: cache}
}

func (ss *SessionManager) CreateSession() (string, error) {
	id := generateSessionID()
	sessionData := &SessionData{}
	sessionBytes, err := json.Marshal(sessionData)
	if err != nil {
		return "", err
	}

	ss.cache.AddSession(id, sessionBytes)
	return id, nil
}

func (ss *SessionManager) UpdateSession(data *SessionData) error {
	//
	return nil
}

func (ss *SessionManager) GetSession(sessionID string) (*SessionData, error) {
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

func (ss *SessionManager) DeleteSession(sessionID string) {
	//
}

func generateSessionID() string {
	return uuid.NewString()
}
