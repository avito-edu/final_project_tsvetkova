package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/casbin/casbin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware(enforcer *casbin.Enforcer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		publicPaths := map[string]bool{
			"/login":    true,
			"/register": true,
		}

		if publicPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Authorization header required"}`, http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"error": "Invalid authorization format"}`, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			http.Error(w, `{"error": "Invalid or expired token: `+err.Error()+`"}`, http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, `{"error": "Invalid token"}`, http.StatusUnauthorized)
			return
		}

		path := r.URL.Path
		method := r.Method

		allowed, err := enforcer.EnforceSafe(claims.Role, path, method)
		if err != nil {
			http.Error(w, `{"error": "Authorization error: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		if !allowed {
			http.Error(w, `{"error": "Access denied"}`, http.StatusForbidden)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "role", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
