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
// Các trường này phải là kiểu giá trị, không phải con trỏ.
type CreateCustomerRequest struct {
	Name    string `json:"name" binding:"required"`
	TaxCode string `json:"tax_code"`
	Address string `json:"address"`
	Email   string `json:"email" binding:"omitempty,email"`
	Phone   string `json:"phone"`
}

// UpdateCustomerRequest cho phép cập nhật một phần thông tin khách hàng.
// Các trường này dùng con trỏ để phân biệt trường nào được gửi lên để cập nhật.
type UpdateCustomerRequest struct {
	Name    *string `json:"name,omitempty"`
	TaxCode *string `json:"tax_code,omitempty"`
	Address *string `json:"address,omitempty"`
	Email   *string `json:"email,omitempty" binding:"omitempty,email"`
	Phone   *string `json:"phone,omitempty"`
}
