package service

import (
	"database/sql"
	"log"

	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/models"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/repository"
	"github.com/hoanggtan02/my-app/vat-simple-backend/pkg/errors"
)

type ProductService interface {
	CreateProduct(req *models.CreateProductRequest, companyID string) (*models.Product, *errors.AppError)
	ListProducts(companyID string) ([]models.Product, *errors.AppError)
	UpdateProduct(productID string, req *models.UpdateProductRequest) *errors.AppError // <-- Thêm hàm mới
}

type productServiceImpl struct {
	productRepo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productServiceImpl{productRepo: repo}
}

func (s *productServiceImpl) CreateProduct(req *models.CreateProductRequest, companyID string) (*models.Product, *errors.AppError) {
	description := sql.NullString{}
	if req.Description != "" {
		description.String = req.Description
		description.Valid = true
	}
	imageUrl := sql.NullString{}
	if req.ImageURL != "" {
		imageUrl.String = req.ImageURL
		imageUrl.Valid = true
	}
	product := &models.Product{
		Name:        req.Name,
		Description: description,
		UnitPrice:   req.UnitPrice,
		ImageURL:    imageUrl,
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

// UpdateProduct (Hàm mới)
func (s *productServiceImpl) UpdateProduct(productID string, req *models.UpdateProductRequest) *errors.AppError {
	if err := s.productRepo.UpdateProduct(productID, req); err != nil {
		log.Printf("Error updating product in repo: %v", err)
		return errors.ErrInternalServerError
	}
	return nil
}
