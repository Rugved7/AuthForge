package middleware

import (
	"net/http"

	contextkeys "github.com/Rugved7/authforge/internal/http/contextKeys"
)

func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	roleSet := make(map[string]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(contextkeys.RoleKey).(string)
			if !ok {
				http.Error(w, "role is missing in context", http.StatusForbidden)
				return
			}

			if _, allowed := roleSet[role]; !allowed {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
