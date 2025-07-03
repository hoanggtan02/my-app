package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/models"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/service"
	"github.com/hoanggtan02/my-app/vat-simple-backend/pkg/errors"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(s service.ProductService) *ProductHandler {
	return &ProductHandler{productService: s}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	companyID, _ := c.Get("companyID")
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(http.StatusBadRequest, "Invalid request body"))
		return
	}
	product, appErr := h.productService.CreateProduct(&req, companyID.(string))
	if appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	companyID, _ := c.Get("companyID")
	products, appErr := h.productService.ListProducts(companyID.(string))
	if appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}
	c.JSON(http.StatusOK, products)
}

// UpdateProduct (Hàm mới)
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID := c.Param("id")
	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(http.StatusBadRequest, "Invalid request body"))
		return
	}

	if appErr := h.productService.UpdateProduct(productID, &req); appErr != nil {
		c.JSON(appErr.StatusCode, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}
