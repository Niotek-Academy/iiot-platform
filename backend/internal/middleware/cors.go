package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS allows the frontend (running on a different origin during local
// dev, e.g. localhost:3000) to call this API. allowedOrigins is a list of
// exact origins to allow, or ["*"] to allow any origin (fine for local
// dev, not recommended once this goes to production with real users).
func CORS(allowedOrigins []string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(allowedOrigins))
	allowAll := false
	for _, o := range allowedOrigins {
		o = strings.TrimSpace(o)
		if o == "*" {
			allowAll = true
		}
		allowed[o] = true
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && (allowAll || allowed[origin]) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// The browser sends an OPTIONS preflight before the real request
		// (that's exactly what's failing in your error) — it just needs a
		// 204 with the headers above, no further processing.
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}