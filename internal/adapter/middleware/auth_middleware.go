package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Vlad-Ali/Movies-service-back/internal/adapter/useridkey"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user"
)

type AuthMiddleware struct {
	tokenService user.TokenService
}

func NewAuthMiddleware(tokenService user.TokenService) *AuthMiddleware {
	return &AuthMiddleware{tokenService}
}

func (am *AuthMiddleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}
		token := parts[1]

		if token != "" {
			userID, err := am.tokenService.ValidateToken(r.Context(), token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), useridkey.UserIDKey{}, userID.ID())
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
