package models

import (
	"database/sql"
	"time"
)

// Product represents a product or service that can be sold.
type Product struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	UnitPrice   float64        `json:"unit_price"`
	ImageURL    sql.NullString `json:"image_url"`
	CompanyID   string         `json:"company_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// CreateProductRequest defines the payload for creating a new product.
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unit_price" binding:"required,gt=0"`
	ImageURL    string  `json:"image_url"`
}

// UpdateProductRequest cho phép cập nhật một phần thông tin sản phẩm
type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	UnitPrice   *float64 `json:"unit_price,omitempty" binding:"omitempty,gt=0"`
	ImageURL    *string  `json:"image_url,omitempty"`
}
