package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/hoanggtan02/my-app/vat-simple-backend/pkg/errors" // Import custom errors package
)

// Validator holds the validator instance
var Validate *validator.Validate

// InitValidator initializes the global validator instance
func InitValidator() {
	Validate = validator.New()
	// Optionally, add custom validation tags or translations here
	// For example:
	// Validate.RegisterValidation("is-custom-valid", ValidateCustomField)
}

// ValidateStruct validates a struct based on its 'binding' tags.
// It translates validation errors into a more user-friendly format (AppError).
func ValidateStruct(s any) *errors.AppError {
	if Validate == nil {
		InitValidator() // Ensure validator is initialized if not already
	}

	if err := Validate.Struct(s); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errorMessages []string
		details := make(map[string]string)

		for _, fieldError := range validationErrors {
			// Example: "email is required" or "password must be at least 8 characters"
			msg := fmt.Sprintf("Field '%s' is invalid due to '%s' validation rule.",
				strings.ToLower(fieldError.Field()), fieldError.Tag())

			// More specific messages for common rules
			switch fieldError.Tag() {
			case "required":
				msg = fmt.Sprintf("Field '%s' is required.", strings.ToLower(fieldError.Field()))
			case "email":
				msg = fmt.Sprintf("Field '%s' must be a valid email address.", strings.ToLower(fieldError.Field()))
			case "min":
				msg = fmt.Sprintf("Field '%s' must be at least %s characters long.", strings.ToLower(fieldError.Field()), fieldError.Param())
			case "max":
				msg = fmt.Sprintf("Field '%s' must not exceed %s characters.", strings.ToLower(fieldError.Field()), fieldError.Param())
			case "gt": // greater than
				msg = fmt.Sprintf("Field '%s' must be greater than %s.", strings.ToLower(fieldError.Field()), fieldError.Param())
			case "oneof": // one of a set of values
				msg = fmt.Sprintf("Field '%s' must be one of: %s.", strings.ToLower(fieldError.Field()), fieldError.Param())
			case "datetime": // specific datetime format
				msg = fmt.Sprintf("Field '%s' must be in format '%s'.", strings.ToLower(fieldError.Field()), fieldError.Param())
			}

			errorMessages = append(errorMessages, msg)
			details[strings.ToLower(fieldError.Field())] = msg // Map field to its error message
		}

		return errors.NewAppError(
			http.StatusBadRequest,
			"Validation failed",
			errors.WithCode("VALIDATION_ERROR"),
			errors.WithDetails(details),
			errors.WithCause(err), // Store the original validation error for logging
		)
	}
	return nil // No validation errors
}
