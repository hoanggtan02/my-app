package errors

import (
	"fmt"
	"net/http"
)

// AppError represents a custom application error with a message and HTTP status code.
type AppError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`                 // Don't marshal status code to JSON, use it for HTTP response
	Code       string `json:"code,omitempty"`    // Custom error code (e.g., "USER_NOT_FOUND", "INVALID_INPUT")
	Details    any    `json:"details,omitempty"` // Optional: additional details about the error
	Err        error  `json:"-"`                 // Original error for logging, not exposed via API
}

// Error implements the error interface for AppError.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewAppError creates a new AppError.
func NewAppError(statusCode int, message string, opts ...AppErrorOption) *AppError {
	err := &AppError{
		StatusCode: statusCode,
		Message:    message,
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}

// AppErrorOption is a functional option for AppError.
type AppErrorOption func(*AppError)

// WithCode sets a custom code for the AppError.
func WithCode(code string) AppErrorOption {
	return func(e *AppError) {
		e.Code = code
	}
}

// WithDetails sets additional details for the AppError.
func WithDetails(details any) AppErrorOption {
	return func(e *AppError) {
		e.Details = details
	}
}

// WithCause sets the underlying error that caused the AppError.
func WithCause(err error) AppErrorOption {
	return func(e *AppError) {
		e.Err = err
	}
}

// Predefined common errors
var (
	ErrNotFound            = NewAppError(http.StatusNotFound, "Resource not found", WithCode("NOT_FOUND"))
	ErrInvalidInput        = NewAppError(http.StatusBadRequest, "Invalid input provided", WithCode("INVALID_INPUT"))
	ErrUnauthorized        = NewAppError(http.StatusUnauthorized, "Authentication required", WithCode("UNAUTHORIZED"))
	ErrForbidden           = NewAppError(http.StatusForbidden, "Access denied", WithCode("FORBIDDEN"))
	ErrInternalServerError = NewAppError(http.StatusInternalServerError, "An unexpected error occurred", WithCode("INTERNAL_SERVER_ERROR"))
	ErrConflict            = NewAppError(http.StatusConflict, "Conflict occurred", WithCode("CONFLICT"))
)
