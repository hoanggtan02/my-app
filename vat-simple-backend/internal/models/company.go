package models

import "time"

// CompanyProfile represents a company's profile in the database
type CompanyProfile struct {
	ID                string    `json:"id" db:"id"`
	UserID            string    `json:"user_id" db:"user_id"`
	Name              string    `json:"name" db:"name"`
	TaxCode           string    `json:"tax_code" db:"tax_code,omitempty"`
	Address           string    `json:"address" db:"address,omitempty"`
	Phone             string    `json:"phone" db:"phone,omitempty"`
	Email             string    `json:"email" db:"email,omitempty"`
	BankAccountNumber string    `json:"bank_account_number" db:"bank_account_number,omitempty"`
	BankName          string    `json:"bank_name" db:"bank_name,omitempty"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// UpdateCompanyProfileRequest represents the request body for updating company profile
type UpdateCompanyProfileRequest struct {
	Name              *string `json:"name,omitempty" binding:"omitempty"`
	TaxCode           *string `json:"tax_code,omitempty" binding:"omitempty"`
	Address           *string `json:"address,omitempty" binding:"omitempty"`
	Phone             *string `json:"phone,omitempty" binding:"omitempty"`
	Email             *string `json:"email,omitempty" binding:"omitempty"`
	BankAccountNumber *string `json:"bank_account_number,omitempty" binding:"omitempty"`
	BankName          *string `json:"bank_name,omitempty" binding:"omitempty"`
}
