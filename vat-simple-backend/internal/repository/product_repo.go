package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/models"
)

type ProductRepository interface {
	CreateProduct(product *models.Product) error
	GetProductByID(productID, companyID string) (*models.Product, error)
	GetAllProducts(companyID string) ([]models.Product, error)
	UpdateProduct(productID string, req *models.UpdateProductRequest) error // <-- Thêm hàm mới vào interface
}

type productRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepositoryImpl{db: db}
}

// CreateProduct (giữ nguyên)
func (r *productRepositoryImpl) CreateProduct(product *models.Product) error {
	product.ID = uuid.New().String()
	query := "INSERT INTO products (id, name, description, unit_price, image_url, company_id) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := r.db.Exec(query, product.ID, product.Name, product.Description, product.UnitPrice, product.ImageURL, product.CompanyID)
	return err
}

// GetProductByID (giữ nguyên)
func (r *productRepositoryImpl) GetProductByID(productID, companyID string) (*models.Product, error) {
	product := &models.Product{}
	query := "SELECT id, name, description, unit_price, image_url FROM products WHERE id = ? AND company_id = ?"
	err := r.db.QueryRow(query, productID, companyID).Scan(&product.ID, &product.Name, &product.Description, &product.UnitPrice, &product.ImageURL)
	return product, err
}

// GetAllProducts (giữ nguyên)
func (r *productRepositoryImpl) GetAllProducts(companyID string) ([]models.Product, error) {
	query := "SELECT id, name, description, unit_price, image_url FROM products WHERE company_id = ? ORDER BY name ASC"
	rows, err := r.db.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.UnitPrice, &product.ImageURL); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

// UpdateProduct (Hàm mới)
func (r *productRepositoryImpl) UpdateProduct(productID string, req *models.UpdateProductRequest) error {
	// Xây dựng câu lệnh UPDATE động
	var setParts []string
	var args []interface{}

	if req.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *req.Name)
	}
	if req.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *req.Description)
	}
	if req.UnitPrice != nil {
		setParts = append(setParts, "unit_price = ?")
		args = append(args, *req.UnitPrice)
	}
	if req.ImageURL != nil {
		setParts = append(setParts, "image_url = ?")
		args = append(args, *req.ImageURL)
	}

	if len(setParts) == 0 {
		return nil // Không có gì để cập nhật
	}

	query := fmt.Sprintf("UPDATE products SET %s WHERE id = ?", strings.Join(setParts, ", "))
	args = append(args, productID)

	_, err := r.db.Exec(query, args...)
	return err
}
