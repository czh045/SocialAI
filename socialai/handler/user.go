package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"socialai/model"
	"socialai/service"

	"github.com/golang-jwt/jwt/v5"
)

var usernamePattern = regexp.MustCompile(`^[A-Za-z0-9_]{3,30}$`)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
		return
	}

	success, err := service.CheckUser(user.Username, user.Password)
	if err != nil {
		fmt.Printf("failed to read user from Elasticsearch: %v\n", err)
		http.Error(w, "Failed to read user from backend", http.StatusInternalServerError)
		return
	}
	if !success {
		http.Error(w, "User does not exist or password is wrong", http.StatusUnauthorized)
		return
	}

	// The user name is stored in the JWT Subject and recovered by auth.go.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   user.Username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(tokenString))
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
		return
	}
	if !usernamePattern.MatchString(user.Username) || len(user.Password) < 3 {
		http.Error(w, "Username must be 3-30 letters, numbers, or underscores; password must be at least 3 characters", http.StatusBadRequest)
		return
	}

	success, err := service.AddUser(&user)
	if err != nil {
		fmt.Printf("failed to save user: %v\n", err)
		http.Error(w, "Failed to save user to backend", http.StatusInternalServerError)
		return
	}
	if !success {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
