package helper

import (
	"net/http"
	"os"
)

func EnableCORS(next http.Handler) http.Handler {
	frontendURL := os.Getenv("FRONTEND_URL")

	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", frontendURL)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}