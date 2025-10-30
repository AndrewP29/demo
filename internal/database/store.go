package database

import (
	"database/sql"

	"demo/internal/models"
)

// Datastore defines the database operations required by the API handlers.
type Datastore interface {
	CreateUser(username, email, hashedPassword string) (int, error)
	GetUserByUsername(username string) (*models.User, string, error)
}

// DBStore is a concrete implementation of the Datastore interface that uses a real SQL database.
type DBStore struct {
	DB *sql.DB
}

// CreateUser inserts a new user into the database.
func (store *DBStore) CreateUser(username, email, hashedPassword string) (int, error) {
	var newUserID int
	err := store.DB.QueryRow("INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
		username, email, hashedPassword).Scan(&newUserID)

	if err != nil {
		return 0, err
	}
	return newUserID, nil
}

// GetUserByUsername retrieves a user and their hashed password from the database.
func (store *DBStore) GetUserByUsername(username string) (*models.User, string, error) {
	var user models.User
	var hashedPassword string

	err := store.DB.QueryRow("SELECT id, username, password_hash FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &hashedPassword)
	if err != nil {
		return nil, "", err
	}
	return &user, hashedPassword, nil
}
