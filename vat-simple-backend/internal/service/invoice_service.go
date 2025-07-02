package service

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/models"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/repository"
	"github.com/hoanggtan02/my-app/vat-simple-backend/pkg/errors"
)

type InvoiceService interface {
	CreateInvoice(req *models.CreateInvoiceRequest, companyID string) (*models.Invoice, *errors.AppError)
	GetInvoice(invoiceID, companyID string) (*models.Invoice, *errors.AppError)
	ListInvoices(companyID string) ([]models.Invoice, *errors.AppError)
}

type invoiceServiceImpl struct {
	invoiceRepo  repository.InvoiceRepository
	customerRepo repository.CustomerRepository
	productRepo  repository.ProductRepository
}

func NewInvoiceService(invoiceRepo repository.InvoiceRepository, customerRepo repository.CustomerRepository, productRepo repository.ProductRepository) InvoiceService {
	return &invoiceServiceImpl{
		invoiceRepo:  invoiceRepo,
		customerRepo: customerRepo,
		productRepo:  productRepo,
	}
}

func (s *invoiceServiceImpl) CreateInvoice(req *models.CreateInvoiceRequest, companyID string) (*models.Invoice, *errors.AppError) {
	// 1. Validate customer existence
	_, err := s.customerRepo.GetCustomerByID(req.CustomerID, companyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewAppError(http.StatusBadRequest, "Customer not found", errors.WithCode("CUSTOMER_NOT_FOUND"))
		}
		return nil, errors.ErrInternalServerError
	}

	// 2. Process invoice items and calculate totals
	var invoiceItems []models.InvoiceItem
	var subtotal float64 = 0

	for _, itemReq := range req.Items {
		product, err := s.productRepo.GetProductByID(itemReq.ProductID, companyID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.NewAppError(http.StatusBadRequest, fmt.Sprintf("Product with ID %s not found", itemReq.ProductID), errors.WithCode("PRODUCT_NOT_FOUND"))
			}
			return nil, errors.ErrInternalServerError
		}

		totalPrice := float64(itemReq.Quantity) * product.UnitPrice
		subtotal += totalPrice

		invoiceItems = append(invoiceItems, models.InvoiceItem{
			ProductID:   itemReq.ProductID,
			Description: product.Name, // Use product name as description
			Quantity:    itemReq.Quantity,
			UnitPrice:   product.UnitPrice,
			TotalPrice:  totalPrice,
		})
	}

	// 3. Calculate tax and total (assuming a flat 10% VAT for simplicity)
	tax := subtotal * 0.10
	total := subtotal + tax

	// Round to 2 decimal places
	subtotal = math.Round(subtotal*100) / 100
	tax = math.Round(tax*100) / 100
	total = math.Round(total*100) / 100

	// 4. Create the main invoice object
	invoice := &models.Invoice{
		CompanyID:     companyID,
		CustomerID:    req.CustomerID,
		InvoiceNumber: req.InvoiceNumber, // In a real app, this should be generated automatically
		IssueDate:     req.IssueDate,
		DueDate:       req.DueDate,
		Subtotal:      subtotal,
		Tax:           tax,
		Total:         total,
		Status:        "draft", // Default status
	}

	// 5. Save to database
	createdInvoice, repoErr := s.invoiceRepo.CreateInvoice(invoice, invoiceItems)
	if repoErr != nil {
		log.Printf("Error creating invoice in repo: %v", repoErr)
		return nil, errors.ErrInternalServerError
	}

	return createdInvoice, nil
}

func (s *invoiceServiceImpl) GetInvoice(invoiceID, companyID string) (*models.Invoice, *errors.AppError) {
	invoice, err := s.invoiceRepo.GetInvoiceByID(invoiceID, companyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewAppError(http.StatusNotFound, "Invoice not found", errors.WithCode("INVOICE_NOT_FOUND"))
		}
		log.Printf("Error getting invoice from repo: %v", err)
		return nil, errors.ErrInternalServerError
	}
	return invoice, nil
}

func (s *invoiceServiceImpl) ListInvoices(companyID string) ([]models.Invoice, *errors.AppError) {
	invoices, err := s.invoiceRepo.GetInvoicesByCompanyID(companyID)
	if err != nil {
		log.Printf("Error listing invoices from repo: %v", err)
		return nil, errors.ErrInternalServerError
	}
	return invoices, nil
}
