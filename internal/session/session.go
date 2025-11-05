package session

import "time"

type Session struct {
	ID     string
	UserID int64
	Expiry time.Time
}

type Store interface {
	Create(userID int64) (string, error)

	Get(sessionID string) (*Session, error)

	Delete(sessionID string) error
}
