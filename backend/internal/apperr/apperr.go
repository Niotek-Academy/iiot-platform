// Package apperr defines a single error type the whole backend uses, so
// every layer (service, handler, middleware) can carry a clear HTTP status
// and machine-readable code alongside the human message.
package apperr

type AppError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string { return e.Message }

func New(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

func NewBadRequest(message string) *AppError  { return New(400, "BAD_REQUEST", message) }
func NewValidation(message string) *AppError  { return New(422, "VALIDATION_ERROR", message) }
func NewUnauthorized(message string) *AppError { return New(401, "UNAUTHORIZED", message) }
func NewForbidden(message string) *AppError    { return New(403, "FORBIDDEN", message) }
func NewNotFound(message string) *AppError     { return New(404, "NOT_FOUND", message) }
func NewConflict(message string) *AppError     { return New(409, "CONFLICT", message) }
func NewInternal(message string) *AppError     { return New(500, "INTERNAL_ERROR", message) }