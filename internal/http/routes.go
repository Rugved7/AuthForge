// Package http provies the http server logic.
package http

import (
	"net/http"

	"github.com/Rugved7/authforge/internal/auth"
)

func NewRouter(authHandler *auth.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/signup", authHandler.Signup)
	mux.HandleFunc("/login", authHandler.Login)

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
}
