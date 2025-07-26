package auth

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/regcomp/gdpr/caching"
)

type SessionData struct {
	//
}

type ISessionStore interface {
	CreateSession() string
	GetSession(string) (*SessionData, error)
	UpdateSession(*SessionData) error
	DeleteSession(string)
}

type SessionStore struct {
	cache caching.IServiceCache
}

func CreateSessionStore(serviceCache caching.IServiceCache) *SessionStore {
	return &SessionStore{cache: serviceCache}
}

func (ss *SessionStore) CreateSession() string {
	id := generateSessionID()
	sessionData := &SessionData{}
	sessionBytes, err := json.Marshal(sessionData)
	if err != nil {
		// TODO:
	}

	ss.cache.SessionAdd(id, sessionBytes)
	return id
}

func (ss *SessionStore) UpdateSession(data *SessionData) error {
	//
	return nil
}

func (ss *SessionStore) GetSession(sessionID string) (*SessionData, error) {
	sessionBytes, err := ss.cache.SessionGet(sessionID)
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
