// Package middleware handles the authorization part
package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Rugved7/authforge/internal/auth"
	"github.com/Rugved7/authforge/internal/cache"
	contextkeys "github.com/Rugved7/authforge/internal/http/contextKeys"
)

func AuthMiddleware(tokenManager *auth.TokenManager, tokenCache cache.Cache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "mising authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]
			cacheKey := "access:" + tokenStr

			// try cache first
			if cached, ok := tokenCache.Get(r.Context(), cacheKey); ok {
				// cached format -> userID:role
				parts := strings.SplitN(cached, ":", 2)
				ctx := context.WithValue(r.Context(), contextkeys.UserIDKey, parts[0])
				ctx = context.WithValue(ctx, contextkeys.RoleKey, parts[1])
				next.ServeHTTP(w, r.WithContext(ctx))
			}

			// fallback to JWT verfication
			_, claims, err := tokenManager.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			// ensure this is accesstoken
			if tokenType, ok := claims["type"].(string); !ok || tokenType != "access" {
				http.Error(w, "invalid token type", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, "invalid token subject", http.StatusUnauthorized)
				return
			}

			role, ok := claims["role"].(string)
			if !ok {
				http.Error(w, "invalid token role", http.StatusUnauthorized)
				return
			}

			// Store in cache
			tokenCache.Set(
				r.Context(),
				cacheKey,
				userID+":"+role,
				5*time.Minute,
			)

			ctx := context.WithValue(r.Context(), contextkeys.UserIDKey, userID)
			ctx = context.WithValue(ctx, contextkeys.RoleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
