package repository

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/models"
)

type InvoiceRepository interface {
	CreateInvoice(invoice *models.Invoice, items []models.InvoiceItem) (*models.Invoice, error)
	GetInvoiceByID(invoiceID, companyID string) (*models.Invoice, error)
	GetInvoicesByCompanyID(companyID string) ([]models.Invoice, error)
}

type invoiceRepositoryImpl struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) InvoiceRepository {
	return &invoiceRepositoryImpl{db: db}
}

// CreateInvoice creates a new invoice and its items in a single DB transaction.
func (r *invoiceRepositoryImpl) CreateInvoice(invoice *models.Invoice, items []models.InvoiceItem) (*models.Invoice, error) {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return nil, err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	// 1. Insert the main invoice record
	invoice.ID = uuid.New().String()
	invoiceQuery := `
		INSERT INTO invoices (id, company_id, customer_id, invoice_number, issue_date, due_date, subtotal, tax, total, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(invoiceQuery,
		invoice.ID, invoice.CompanyID, invoice.CustomerID, invoice.InvoiceNumber, invoice.IssueDate, invoice.DueDate, invoice.Subtotal, invoice.Tax, invoice.Total, invoice.Status,
	)
	if err != nil {
		log.Printf("Failed to insert invoice: %v", err)
		return nil, err
	}

	// 2. Prepare statement for inserting invoice items
	itemQuery := `
		INSERT INTO invoice_items (id, invoice_id, product_id, description, quantity, unit_price, total_price)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	stmt, err := tx.Prepare(itemQuery)
	if err != nil {
		log.Printf("Failed to prepare statement for invoice items: %v", err)
		return nil, err
	}
	defer stmt.Close()

	// 3. Insert each invoice item
	for _, item := range items {
		_, err := stmt.Exec(uuid.New().String(), invoice.ID, item.ProductID, item.Description, item.Quantity, item.UnitPrice, item.TotalPrice)
		if err != nil {
			log.Printf("Failed to insert invoice item: %v", err)
			return nil, err
		}
	}

	// If all queries were successful, commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return nil, err
	}

	return invoice, nil
}

// GetInvoiceByID retrieves a single invoice and its associated items.
func (r *invoiceRepositoryImpl) GetInvoiceByID(invoiceID, companyID string) (*models.Invoice, error) {
	invoice := &models.Invoice{}
	query := `
		SELECT i.id, i.company_id, i.customer_id, i.invoice_number, i.issue_date, i.due_date, i.subtotal, i.tax, i.total, i.status, c.name as customer_name
		FROM invoices i
		JOIN customers c ON i.customer_id = c.id
		WHERE i.id = ? AND i.company_id = ?`

	err := r.db.QueryRow(query, invoiceID, companyID).Scan(
		&invoice.ID, &invoice.CompanyID, &invoice.CustomerID, &invoice.InvoiceNumber, &invoice.IssueDate, &invoice.DueDate, &invoice.Subtotal, &invoice.Tax, &invoice.Total, &invoice.Status, &invoice.CustomerName,
	)
	if err != nil {
		return nil, err // Returns sql.ErrNoRows if not found
	}

	// Now fetch the items for the invoice
	itemsQuery := `SELECT id, product_id, description, quantity, unit_price, total_price FROM invoice_items WHERE invoice_id = ?`
	rows, err := r.db.Query(itemsQuery, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.InvoiceItem
	for rows.Next() {
		var item models.InvoiceItem
		if err := rows.Scan(&item.ID, &item.ProductID, &item.Description, &item.Quantity, &item.UnitPrice, &item.TotalPrice); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	invoice.Items = items

	return invoice, nil
}

// GetInvoicesByCompanyID retrieves all invoices for a given company.
func (r *invoiceRepositoryImpl) GetInvoicesByCompanyID(companyID string) ([]models.Invoice, error) {
	query := `
		SELECT i.id, i.invoice_number, i.issue_date, i.total, i.status, c.name as customer_name
		FROM invoices i
		JOIN customers c ON i.customer_id = c.id
		WHERE i.company_id = ?
		ORDER BY i.issue_date DESC`

	rows, err := r.db.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		var invoice models.Invoice
		if err := rows.Scan(&invoice.ID, &invoice.InvoiceNumber, &invoice.IssueDate, &invoice.Total, &invoice.Status, &invoice.CustomerName); err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}
	return invoices, nil
}
