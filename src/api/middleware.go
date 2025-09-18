package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type contextKey string

const UserIDContextKey contextKey = "userID"

type authResponse struct {
	IsValid bool   `json:"is_valid"`
	UserID  string `json:"user_id,omitempty"`
}

func (h *Handler) JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Missing or malformed token", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		requestBody, _ := json.Marshal(map[string]string{"token": token})
		authReq, _ := http.NewRequest(http.MethodPost, h.cfg.AuthServiceURL, bytes.NewBuffer(requestBody))
		authReq.Header.Set("Content-Type", "application/json")
		authReq.Header.Set("X-Internal-Api-Key", h.cfg.InternalAPIKey)

		client := &http.Client{}
		res, err := client.Do(authReq)
		if err != nil {
			http.Error(w, "Internal Server Error: Could not reach auth service", http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		var authRes authResponse
		if err := json.NewDecoder(res.Body).Decode(&authRes); err != nil {
			http.Error(w, "Internal Server Error: Could not decode auth response", http.StatusInternalServerError)
			return
		}

		if !authRes.IsValid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDContextKey, authRes.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) APIKeyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		providedKey := r.Header.Get("X-Internal-Api-Key")
		if providedKey == "" || providedKey != h.cfg.InternalAPIKey {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
