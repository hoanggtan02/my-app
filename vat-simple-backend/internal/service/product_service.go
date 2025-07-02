package service

import (
	"log"

	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/models"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/repository"
	"github.com/hoanggtan02/my-app/vat-simple-backend/pkg/errors"
)

type ProductService interface {
	CreateProduct(req *models.CreateProductRequest, companyID string) (*models.Product, *errors.AppError)
	ListProducts(companyID string) ([]models.Product, *errors.AppError)
}

type productServiceImpl struct {
	productRepo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productServiceImpl{productRepo: repo}
}

func (s *productServiceImpl) CreateProduct(req *models.CreateProductRequest, companyID string) (*models.Product, *errors.AppError) {
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		UnitPrice:   req.UnitPrice, // <-- Lỗi xảy ra ở đây, code này đã đúng
		CompanyID:   companyID,
	}

	if err := s.productRepo.CreateProduct(product); err != nil {
		log.Printf("Error creating product in repo: %v", err)
		return nil, errors.ErrInternalServerError
	}
	return product, nil
}

func (s *productServiceImpl) ListProducts(companyID string) ([]models.Product, *errors.AppError) {
	products, err := s.productRepo.GetAllProducts(companyID)
	if err != nil {
		log.Printf("Error listing products from repo: %v", err)
		return nil, errors.ErrInternalServerError
	}
	return products, nil
}
