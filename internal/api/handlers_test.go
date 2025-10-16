package api

import (
	"bytes"
	"demo/internal/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// mockStore is a mock implementation of the Datastore interface for testing.
type mockStore struct {
	users map[string]*models.User
	hashedPasswords map[string]string
}

// GetUserByUsername simulates fetching a user from the mock store.
func (m *mockStore) GetUserByUsername(username string) (*models.User, string, error) {
	if user, ok := m.users[username]; ok {
		return user, m.hashedPasswords[username], nil
	}
	return nil, "", errors.New("user not found")
}

// CreateUser is a mock implementation and is not needed for this login test.
func (m *mockStore) CreateUser(username, email, hashedPassword string) (int, error) {
	// Not needed for this test, but required to satisfy the interface
	return 0, nil
}

func TestLoginHandler(t *testing.T) {
	// Create a dummy hashed password for our test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Setup our mock store with a test user
	mock := &mockStore{
		users: map[string]*models.User{
			"testuser": {ID: 1, Username: "testuser"},
		},
		hashedPasswords: map[string]string{
			"testuser": string(hashedPassword),
		},
	}

	// Create a server instance with our mock store
	server := &Server{Store: mock}

	// Define our table of test cases
	testCases := []struct {
		name               string
		payload            map[string]string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "Successful Login",
			payload:            map[string]string{"username": "testuser", "password": "password123"},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"Login successful","userId":1}`,
		},
		{
			name:               "User Not Found",
			payload:            map[string]string{"username": "nouser", "password": "password123"},
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"error":"Invalid username or password"}`,
		},
		{
			name:               "Incorrect Password",
			payload:            map[string]string{"username": "testuser", "password": "wrongpassword"},
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"error":"Invalid username or password"}`,
		},
		{
			name:               "Malformed JSON",
			payload:            nil, // We will send a malformed body for this case
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"Invalid request body"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var requestBody []byte
			if tc.name == "Malformed JSON" {
				requestBody = []byte(`{"username":"testuser"`) // Incomplete JSON
			} else {
				requestBody, _ = json.Marshal(tc.payload)
			}

			req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(requestBody))
			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(server.LoginHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatusCode)
			}

			// Trim newline characters from the response body for consistent comparison
			gotBody := bytes.TrimSpace(rr.Body.Bytes())
			wantBody := []byte(tc.expectedBody)

			if !bytes.Equal(gotBody, wantBody) {
				t.Errorf("handler returned unexpected body: got %s want %s", gotBody, wantBody)
			}
		})
	}
}
