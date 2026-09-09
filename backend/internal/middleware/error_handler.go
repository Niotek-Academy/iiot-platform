package middleware  // consistent response, and a panic anywhere is caught the same way.

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/apperr"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/response"
)


func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var appErr *apperr.AppError
		if errors.As(err, &appErr) {
			response.Fail(c, appErr)
			return
		}

		log.Printf("unhandled error: %v", err)
		response.FailGeneric(c)
	}
}

// Recovery catches panics anywhere in a handler and turns them into the
// same consistent JSON error shape instead of a raw 500 / connection reset.
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("panic recovered: %v", recovered)
		response.FailGeneric(c)
		c.Abort()
	})
}