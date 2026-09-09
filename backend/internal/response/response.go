package response // Package response gives every endpoint in the API the same JSON shape,

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/apperr"
)

type Envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Success writes {"success": true, "data": ...}.
func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, Envelope{Success: true, Data: data})
}

// Fail writes {"success": false, "error": {...}} using the AppError's own
// HTTP status.
func Fail(c *gin.Context, err *apperr.AppError) {
	c.JSON(err.Status, Envelope{Success: false, Error: &ErrorBody{Code: err.Code, Message: err.Message}})
}

// FailGeneric is the last-resort fallback for an error nobody classified.
func FailGeneric(c *gin.Context) {
	Fail(c, apperr.NewInternal("internal server error"))
	_ = http.StatusInternalServerError // kept only to document the mapping above
}