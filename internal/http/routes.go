package http

import (
	"net/http"

	"github.com/Rugved7/authforge/internal/auth"
	"github.com/Rugved7/authforge/internal/cache"
	"github.com/Rugved7/authforge/internal/http/middleware"
)

func NewRouter(
	authHandler *auth.Handler,
	tokenManager *auth.TokenManager,
	tokenCache cache.Cache,
) http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/signup", authHandler.Signup)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/refresh", authHandler.Refresh)

	// Auth middleware (with cache)
	authMW := middleware.AuthMiddleware(tokenManager, tokenCache)

	// Protected routes
	mux.Handle("/me", authMW(http.HandlerFunc(authHandler.Me)))

	// Admin-only route
	adminOnly := middleware.RequireRole("admin")(
		authMW(http.HandlerFunc(authHandler.AdminPing)),
	)
	mux.Handle("/admin/ping", adminOnly)

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
