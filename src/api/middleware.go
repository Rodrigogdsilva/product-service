package api

import (
	"net/http"
)

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
