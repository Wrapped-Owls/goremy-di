package middleware

import (
	"net/http"
	"strings"

	"github.com/wrapped-owls/goremy-di/examples/context_jwt_user/internal/middleware/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		scheme, token, _ := strings.Cut(authHeader, " ")
		if scheme != "Bearer" || strings.TrimSpace(token) == "" {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ParseClaimsPayload(token)
		if err != nil {
			http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := auth.ContextWithClaims(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
