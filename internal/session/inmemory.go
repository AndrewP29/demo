package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

// MemoryStore is an in-memory implementation of the session.Store interface.
// It is safe for concurrent use.
type MemoryStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

// NewMemoryStore creates and returns a new MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		sessions: make(map[string]*Session),
	}
}

// Create creates a new session for the given userID, stores it, and returns the session ID.
func (s *MemoryStore) Create(userID int64) (string, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[sessionID] = &Session{
		ID:     sessionID,
		UserID: userID,
		Expiry: time.Now().Add(24 * time.Hour), // Session expires in 24 hours
	}

	return sessionID, nil
}

// Get retrieves a session by its ID. It returns an error if the session is not found or is expired.
func (s *MemoryStore) Get(sessionID string) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return nil, errors.New("session not found")
	}

	if time.Now().After(session.Expiry) {
		return nil, errors.New("session expired")
	}

	return session, nil
}

// Delete removes a session by its ID.
func (s *MemoryStore) Delete(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, sessionID)
	return nil
}

// generateSessionID creates a new, cryptographically secure, random session ID.
func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
