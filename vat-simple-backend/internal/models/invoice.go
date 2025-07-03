package models

import "time"

// Invoice represents a single invoice.
type Invoice struct {
	ID            string        `json:"id"`
	CompanyID     string        `json:"company_id"`
	CustomerID    string        `json:"customer_id"`
	CustomerName  string        `json:"customer_name,omitempty"`
	InvoiceNumber string        `json:"invoice_number"`
	IssueDate     time.Time     `json:"issue_date"`
	DueDate       time.Time     `json:"due_date"`
	Subtotal      float64       `json:"subtotal"`
	Tax           float64       `json:"tax"`
	Total         float64       `json:"total"`
	Status        string        `json:"status"` // e.g., "draft", "sent", "paid"
	Items         []InvoiceItem `json:"items,omitempty"`
}

// InvoiceItem represents a single line item on an invoice.
type InvoiceItem struct {
	ID          string  `json:"id"`
	InvoiceID   string  `json:"-"` // Foreign key, hidden from JSON response
	ProductID   string  `json:"product_id"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}

// CreateInvoiceRequest defines the payload for creating a new invoice.
type CreateInvoiceRequest struct {
	Items []struct {
		ProductID string `json:"product_id" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required,gt=0"`
	} `json:"items" binding:"required,min=1"`
}

// CreateInvoiceItemRequest is a sub-struct for invoice creation.
type CreateInvoiceItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"` // <-- SỬA TỪ FLOAT64 SANG INT
}
