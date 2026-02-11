package middleware

import (
	"context"
	"net/http"

	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

type TokenClaims struct {
	UserID string
	Role   string
}

type ValidateTokenFunc func(token string) (*TokenClaims, error)

func Auth(validate ValidateTokenFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				lib.RespondError(w, http.StatusUnauthorized, "missing access_token cookie")
				return
			}

			claims, err := validate(cookie.Value)
			if err != nil {
				lib.RespondError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(UserIDKey).(string); ok {
		return v
	}
	return ""
}

func GetRole(ctx context.Context) string {
	if v, ok := ctx.Value(RoleKey).(string); ok {
		return v
	}
	return ""
}

// NewValidateTokenFunc cria um ValidateTokenFunc a partir de uma função de extração
// Uso: middleware.Auth(middleware.NewValidateTokenFunc(extractClaims))
func NewValidateTokenFunc(extract func(token string) (userID string, role string, err error)) ValidateTokenFunc {
	return func(token string) (*TokenClaims, error) {
		userID, role, err := extract(token)
		if err != nil {
			return nil, err
		}
		return &TokenClaims{
			UserID: userID,
			Role:   role,
		}, nil
	}
}
