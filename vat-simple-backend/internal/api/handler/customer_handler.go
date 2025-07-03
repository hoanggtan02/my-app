package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/models"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/service"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/utils"
	"github.com/hoanggtan02/my-app/vat-simple-backend/pkg/errors"
)

type CustomerHandler struct {
	customerService service.CustomerService
}

func NewCustomerHandler(s service.CustomerService) *CustomerHandler {
	return &CustomerHandler{customerService: s}
}

// CreateCustomer handles the creation of a new customer.
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	companyID, _ := c.Get("companyID")

	var req models.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if appErr := utils.ValidateStruct(&req); appErr != nil {
			c.JSON(appErr.StatusCode, appErr)
			return
		}
		c.JSON(http.StatusBadRequest, errors.NewAppError(http.StatusBadRequest, "Invalid request body", errors.WithCause(err)))
		return
	}

	customer, appErr := h.customerService.CreateCustomer(&req, companyID.(string))
	if appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusCreated, customer)
}

// GetCustomer handles retrieving a single customer.
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	companyID, _ := c.Get("companyID")
	customerID := c.Param("id")

	customer, appErr := h.customerService.GetCustomer(customerID, companyID.(string))
	if appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusOK, customer)
}

// ListCustomers handles listing all customers for a company.
func (h *CustomerHandler) ListCustomers(c *gin.Context) {
	companyID, _ := c.Get("companyID")
	log.Printf("[DEBUG] Đang lấy khách hàng cho company_id: %s", companyID.(string))

	customers, appErr := h.customerService.ListCustomers(companyID.(string))
	if appErr != nil {
		log.Printf("[ERROR] Service trả về lỗi: %v", appErr)
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	// Dòng log quan trọng nhất: In ra dữ liệu ngay trước khi gửi về client
	log.Printf("[DEBUG] Dữ liệu chuẩn bị gửi về client (%d items): %+v", len(customers), customers)

	c.JSON(http.StatusOK, customers)
}

// UpdateCustomer handles updating an existing customer.
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	companyID, _ := c.Get("companyID")
	customerID := c.Param("id")

	var req models.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(http.StatusBadRequest, "Invalid request body"))
		return
	}

	// Sửa lại dòng này: chỉ nhận 1 giá trị trả về là lỗi
	if appErr := h.customerService.UpdateCustomer(customerID, companyID.(string), &req); appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully"})
}

// DeleteCustomer handles deleting a customer.
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	companyID, _ := c.Get("companyID")
	customerID := c.Param("id")

	if appErr := h.customerService.DeleteCustomer(customerID, companyID.(string)); appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusNoContent, nil) // 204 No Content for successful deletion
}
