package handlers 

// bind.go centralizes request validation so every request can be validated in the same way. 

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/apperr"
)


func BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(apperr.NewValidation(formatValidationError(err)))
		return false
	}
	return true
}

func formatValidationError(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		msgs := make([]string, 0, len(ve))
		for _, fe := range ve {
			msgs = append(msgs, fmt.Sprintf("%s: failed on '%s'", fe.Field(), fe.Tag()))
		}
		return strings.Join(msgs, "; ")
	}
	return err.Error()
}