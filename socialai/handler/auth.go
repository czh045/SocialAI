package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const usernameContextKey contextKey = "username"

func withJWTAuth(next http.Handler, signingKey []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, "Missing bearer token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return signingKey, nil
		})
		if err != nil || !token.Valid || claims.Subject == "" {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), usernameContextKey, claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func currentUsername(r *http.Request) (string, bool) {
	username, ok := r.Context().Value(usernameContextKey).(string)
	return username, ok
}
