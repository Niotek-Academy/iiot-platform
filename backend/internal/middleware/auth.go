package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/apperr"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/utils/jwtutil"
)

// Auth requires a valid "Authorization: Bearer <token>" header. On success
// it stores user_id, username, and role in the Gin context for downstream
// handlers and RequireRole to use.
func Auth(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := parseBearer(c, secret)
		if !ok {
			c.Error(apperr.NewUnauthorized("missing, malformed, or expired token"))
			c.Abort()
			return
		}
		setClaims(c, claims)
		c.Next()
	}
}

// OptionalAuth parses the token if present but never rejects the request
// when it's missing or invalid — used on /auth/register, which must stay
// reachable with no token at all for the very first user (bootstrap).
func OptionalAuth(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		if claims, ok := parseBearer(c, secret); ok {
			setClaims(c, claims)
		}
		c.Next()
	}
}

func parseBearer(c *gin.Context, secret []byte) (*jwtutil.Claims, bool) {
	header := c.GetHeader("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return nil, false
	}
	tokenStr := strings.TrimPrefix(header, "Bearer ")
	claims, err := jwtutil.Parse(secret, tokenStr)
	if err != nil {
		return nil, false
	}
	return claims, true
}

func setClaims(c *gin.Context, claims *jwtutil.Claims) {
	c.Set("user_id", claims.UserID)
	c.Set("username", claims.Username)
	c.Set("role", claims.Role)
}