package middleware

import (
	"context"
	"net/http"
	"strings"

	"qubix/internal/model"
	"qubix/internal/service"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

func AuthMiddleware(jwtService *service.JWTService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, model.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			http.Error(w, model.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		userID, role, err := jwtService.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, model.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserRoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
