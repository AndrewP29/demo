package database

import (
	"database/sql"
	"os"
	"testing"
)

// testDB is a helper function to create a connection to the test database.
func testDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	// Check for environment variables
	if os.Getenv("DB_USER") == "" || os.Getenv("DB_PASSWORD") == "" || os.Getenv("DB_NAME") == "" {
		t.Skip("Skipping database integration tests: DB_USER, DB_PASSWORD, or DB_NAME not set")
	}

	db, err := NewDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Teardown function to clean up the database after the test
	teardown := func() {
		// Truncate tables to reset state
		_ = db.QueryRow("TRUNCATE TABLE users CASCADE")
		db.Close()
	}

	return db, teardown
}

func TestDBStore_CreateUser(t *testing.T) {
	db, teardown := testDB(t)
	defer teardown()

	store := &DBStore{DB: db}

	username := "testuser"
	email := "test@example.com"
	hashedPassword := "hashedpassword"

	id, err := store.CreateUser(username, email, hashedPassword)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if id == 0 {
		t.Errorf("expected non-zero ID, got %d", id)
	}
}

func TestDBStore_GetUserByUsername(t *testing.T) {
	db, teardown := testDB(t)
	defer teardown()

	store := &DBStore{DB: db}

	// First, create a user to fetch
	wantUsername := "testgetuser"
	wantEmail := "get@example.com"
	wantHashedPassword := "$2a$10$...."

	_, err := db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)",
		wantUsername, wantEmail, wantHashedPassword)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	// Now, test the GetUserByUsername method
	user, hashedPassword, err := store.GetUserByUsername(wantUsername)
	if err != nil {
		t.Fatalf("GetUserByUsername failed: %v", err)
	}

	if user.Username != wantUsername {
		t.Errorf("got username %q, want %q", user.Username, wantUsername)
	}

	if hashedPassword != wantHashedPassword {
		t.Errorf("got hashed password %q, want %q", hashedPassword, wantHashedPassword)
	}
}
