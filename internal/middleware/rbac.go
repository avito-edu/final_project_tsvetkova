package middleware

import (
	"os"
	"strings"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		publicPaths := map[string]bool{
			"/login":          true,
			"/register":       true,
			"/register-special": true,
			"/health":         true,
		}

		if publicPaths[c.Request.URL.Path] {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": "Invalid authorization format"})
			c.Abort()
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
			c.JSON(401, gin.H{"error": "Invalid or expired token: " + err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		path := c.Request.URL.Path
		method := c.Request.Method

		allowed, err := enforcer.EnforceSafe(claims.Role, path, method)
		if err != nil {
			c.JSON(500, gin.H{"error": "Authorization error: " + err.Error()})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(403, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}