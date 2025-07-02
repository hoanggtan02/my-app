package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/models"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/service"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/utils"
	"github.com/hoanggtan02/my-app/vat-simple-backend/pkg/errors"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(s service.ProductService) *ProductHandler {
	return &ProductHandler{productService: s}
}

// CreateProduct handles creating a new product.
func (h *ProductHandler) CreateProduct(c *gin.Context) {
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
func (h *ProductHandler) ListProducts(c *gin.Context) {
	companyID, _ := c.Get("companyID")

	products, appErr := h.productService.ListProducts(companyID.(string))
	if appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusOK, products)
}
