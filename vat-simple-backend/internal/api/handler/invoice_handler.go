package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/models"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/service"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/utils"
	"github.com/hoanggtan02/my-app/vat-simple-backend/pkg/errors"
)

type InvoiceHandler struct {
	invoiceService service.InvoiceService
	productService service.ProductService // Also need product service to list products
}

func NewInvoiceHandler(is service.InvoiceService, ps service.ProductService) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceService: is,
		productService: ps,
	}
}

// CreateInvoice handles the creation of a new invoice.
func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	companyID, _ := c.Get("companyID")

	var req models.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if appErr := utils.ValidateStruct(&req); appErr != nil {
			c.JSON(appErr.StatusCode, appErr)
			return
		}
		c.JSON(http.StatusBadRequest, errors.NewAppError(http.StatusBadRequest, "Invalid request body", errors.WithCause(err)))
		return
	}

	invoice, appErr := h.invoiceService.CreateInvoice(&req, companyID.(string))
	if appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusCreated, invoice)
}

// GetInvoice handles retrieving a single invoice.
func (h *InvoiceHandler) GetInvoice(c *gin.Context) {
	companyID, _ := c.Get("companyID")
	invoiceID := c.Param("id")

	invoice, appErr := h.invoiceService.GetInvoice(invoiceID, companyID.(string))
	if appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusOK, invoice)
}

// ListInvoices handles listing all invoices for a company.
func (h *InvoiceHandler) ListInvoices(c *gin.Context) {
	companyID, _ := c.Get("companyID")

	invoices, appErr := h.invoiceService.ListInvoices(companyID.(string))
	if appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusOK, invoices)
}

// --- Product specific handlers needed for invoice creation ---

// CreateProduct handles creating a new product.
func (h *InvoiceHandler) CreateProduct(c *gin.Context) {
	companyID, _ := c.Get("companyID")

	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if appErr := utils.ValidateStruct(&req); appErr != nil {
			c.JSON(appErr.StatusCode, appErr)
			return
		}
		c.JSON(http.StatusBadRequest, errors.NewAppError(http.StatusBadRequest, "Invalid request body", errors.WithCause(err)))
		return
	}

	product, appErr := h.productService.CreateProduct(&req, companyID.(string))
	if appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusCreated, product)
}

// ListProducts handles listing all products.
func (h *InvoiceHandler) ListProducts(c *gin.Context) {
	companyID, _ := c.Get("companyID")

	products, appErr := h.productService.ListProducts(companyID.(string))
	if appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusOK, products)
}
