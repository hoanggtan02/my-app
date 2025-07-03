package service

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

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
	// --- Tự động tìm hoặc tạo "Khách vãng lai" ---
	customer, err := s.customerRepo.FindCustomerByName("Khách vãng lai", companyID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Nếu chưa có, tạo mới
			walkInCustomer := &models.Customer{
				Name:      "Khách vãng lai",
				CompanyID: companyID,
			}
			err = s.customerRepo.CreateCustomer(walkInCustomer)
			if err != nil {
				log.Printf("Error creating walk-in customer: %v", err)
				return nil, errors.ErrInternalServerError
			}
			customer = walkInCustomer
		} else {
			return nil, errors.ErrInternalServerError
		}
	}

	// --- Xử lý sản phẩm và tính toán (giữ nguyên logic cũ) ---
	var invoiceItems []models.InvoiceItem
	var subtotal float64 = 0
	for _, itemReq := range req.Items {
		product, err := s.productRepo.GetProductByID(itemReq.ProductID, companyID)
		if err != nil {
			return nil, errors.NewAppError(http.StatusBadRequest, fmt.Sprintf("Product with ID %s not found", itemReq.ProductID), errors.WithCode("PRODUCT_NOT_FOUND"))
		}
		totalPrice := float64(itemReq.Quantity) * product.UnitPrice
		subtotal += totalPrice
		invoiceItems = append(invoiceItems, models.InvoiceItem{
			ProductID:   itemReq.ProductID,
			Description: product.Name,
			Quantity:    itemReq.Quantity,
			UnitPrice:   product.UnitPrice,
			TotalPrice:  totalPrice,
		})
	}
	tax := subtotal * 0.10
	total := subtotal + tax

	// --- Tạo đối tượng hóa đơn với các thông tin tự động ---
	invoice := &models.Invoice{
		CompanyID:     companyID,
		CustomerID:    customer.ID,                              // <-- Dùng ID của khách vãng lai
		InvoiceNumber: fmt.Sprintf("POS-%d", time.Now().Unix()), // <-- Tự tạo số hóa đơn
		IssueDate:     time.Now(),                               // <-- Lấy ngày hiện tại
		DueDate:       time.Now(),                               // <-- Lấy ngày hiện tại
		Subtotal:      math.Round(subtotal*100) / 100,
		Tax:           math.Round(tax*100) / 100,
		Total:         math.Round(total*100) / 100,
		Status:        "paid", // Trạng thái "đã thanh toán"
	}

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
