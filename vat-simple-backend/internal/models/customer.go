package models

import "time"

// Customer represents a customer of a company.
type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	TaxCode   string    `json:"tax_code"`
	Address   string    `json:"address"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CompanyID string    `json:"company_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateCustomerRequest defines the payload for creating a new customer.
type CreateCustomerRequest struct {
	Name    string `json:"name" binding:"required"`
	TaxCode string `json:"tax_code"`
	Address string `json:"address"`
	Email   string `json:"email" binding:"omitempty,email"`
	Phone   string `json:"phone"`
}

// UpdateCustomerRequest defines the payload for updating a customer.
type UpdateCustomerRequest struct {
	Name    string `json:"name" binding:"required"`
	TaxCode string `json:"tax_code"`
	Address string `json:"address"`
	Email   string `json:"email" binding:"omitempty,email"`
	Phone   string `json:"phone"`
}
