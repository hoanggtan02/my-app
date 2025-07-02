package handler

import (
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

	customers, appErr := h.customerService.ListCustomers(companyID.(string))
	if appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusOK, customers)
}

// UpdateCustomer handles updating an existing customer.
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	companyID, _ := c.Get("companyID")
	customerID := c.Param("id")

	var req models.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if appErr := utils.ValidateStruct(&req); appErr != nil {
			c.JSON(appErr.StatusCode, appErr)
			return
		}
		c.JSON(http.StatusBadRequest, errors.NewAppError(http.StatusBadRequest, "Invalid request body", errors.WithCause(err)))
		return
	}

	updatedCustomer, appErr := h.customerService.UpdateCustomer(customerID, &req, companyID.(string))
	if appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusOK, updatedCustomer)
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
