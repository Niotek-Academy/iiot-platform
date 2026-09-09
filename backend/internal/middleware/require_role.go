package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/apperr"
)

// RequireRole must run after Auth(). It rejects the request with 403 if the
// authenticated user's role isn't in the allowed list.
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}

	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		role, _ := roleVal.(string)

		if !exists || !allowed[role] {
			c.Error(apperr.NewForbidden("you don't have permission to perform this action"))
			c.Abort()
			return
		}
		c.Next()
	}
}