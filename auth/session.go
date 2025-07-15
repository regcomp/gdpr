package auth

import (
	"fmt"

	"github.com/google/uuid"
)

type SessionData struct {
	//
}

type ISessionStore interface {
	CreateSession() string
	GetSession(string) (*SessionData, error)
	UpdateSession(SessionData) error
	DeleteSession(string)
}

type SessionStore struct {
	sessionIDToSessionData map[string]*SessionData
}

func CreateSessionStore() *SessionStore {
	return &SessionStore{
		sessionIDToSessionData: make(map[string]*SessionData),
	}
}

func (ss *SessionStore) CreateSession() string {
	id := generateSessionID()
	data := &SessionData{}
	ss.sessionIDToSessionData[id] = data
	return id
}

func (ss *SessionStore) UpdateSession(data SessionData) error {
	//
	return nil
}

func (ss *SessionStore) GetSession(sessionID string) (*SessionData, error) {
	sessionData, ok := ss.sessionIDToSessionData[sessionID]
	if !ok {
		return nil, fmt.Errorf("could not find session")
	}

	return sessionData, nil
}

func (ss *SessionStore) DeleteSession(sessionID string) {
	//
}

func generateSessionID() string {
	return uuid.NewString()
}
