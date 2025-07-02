package models

import "time"

// User represents a user account in the system.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Hide password hash from JSON responses
	CompanyID    string    `json:"company_id"`
	CompanyName  string    `json:"company_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RegisterRequest defines the payload for user registration.
type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	CompanyName string `json:"company_name" binding:"required"`
}

// LoginRequest defines the payload for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserProfileRequest defines the payload for updating a user's company info.
type UpdateUserProfileRequest struct {
	CompanyName string `json:"company_name" binding:"required"`
}
