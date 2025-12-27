// Package http provies the http server logic.
package http

import (
	"net/http"

	"github.com/Rugved7/authforge/internal/auth"
	"github.com/Rugved7/authforge/internal/http/middleware"
)

func NewRouter(authHandler *auth.Handler, tokenManager *auth.TokenManager) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/signup", authHandler.Signup)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/refresh", authHandler.Refresh)

	// Protected routes
	protected := middleware.AuthMiddleware(tokenManager)(
		http.HandlerFunc(authHandler.Me),
	)
	mux.Handle("/me", protected)

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
}
