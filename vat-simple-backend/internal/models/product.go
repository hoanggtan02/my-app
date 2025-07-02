package models

import "time"

// Product represents a product or service in the database
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UnitPrice   float64   `json:"unit_price"`
	CompanyID   string    `json:"company_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductRequest defines the payload for creating a new product.
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unit_price" binding:"required,gt=0"`
}

// UpdateProductRequest represents the request body for updating an existing product/service
type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" binding:"omitempty"`
	Unit        *string  `json:"unit,omitempty" binding:"omitempty"`
	Price       *float64 `json:"price,omitempty" binding:"omitempty,gt=0"`
	VATRate     *float64 `json:"vat_rate,omitempty" binding:"omitempty,min=0,max=1"`
	Description *string  `json:"description,omitempty" binding:"omitempty"`
}
