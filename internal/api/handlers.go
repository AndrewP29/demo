package api

import (
	"demo/internal/database"
	"demo/internal/models"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Server holds the dependencies for the API handlers, like the datastore.
type Server struct {
	Store database.Datastore
}

func (s *Server) SignupHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming JSON request into a User struct
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error": "Failed to hash password"}`, http.StatusInternalServerError)
		return
	}

	// Use the datastore to create the user
	newUserID, err := s.Store.CreateUser(user.Username, user.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to create user: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Printf("New user created: %s (ID: %d)\n", user.Username, newUserID)

	// Send back a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created successfully",
		"userId":  newUserID,
	})
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming JSON request
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Get user from database using the datastore
	user, hashedPassword, err := s.Store.GetUserByUsername(creds.Username)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid username or password"})
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid username or password"})
		return
	}

	fmt.Printf("User logged in: %s (ID: %d)\n", user.Username, user.ID)

	// Send back a success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"userId":  user.ID,
	})
}