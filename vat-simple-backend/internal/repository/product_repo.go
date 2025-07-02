package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/models"
)

type ProductRepository interface {
	CreateProduct(product *models.Product) error
	GetProductByID(productID, companyID string) (*models.Product, error)
	GetAllProducts(companyID string) ([]models.Product, error)
}

type productRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepositoryImpl{db: db}
}

func (r *productRepositoryImpl) CreateProduct(product *models.Product) error {
	product.ID = uuid.New().String()
	query := "INSERT INTO products (id, name, description, unit_price, company_id) VALUES (?, ?, ?, ?, ?)"
	_, err := r.db.Exec(query, product.ID, product.Name, product.Description, product.UnitPrice, product.CompanyID)
	return err
}

func (r *productRepositoryImpl) GetProductByID(productID, companyID string) (*models.Product, error) {
	product := &models.Product{}
	query := "SELECT id, name, description, unit_price FROM products WHERE id = ? AND company_id = ?"
	err := r.db.QueryRow(query, productID, companyID).Scan(&product.ID, &product.Name, &product.Description, &product.UnitPrice)
	return product, err
}

func (r *productRepositoryImpl) GetAllProducts(companyID string) ([]models.Product, error) {
	query := "SELECT id, name, description, unit_price FROM products WHERE company_id = ? ORDER BY name ASC"
	rows, err := r.db.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.UnitPrice); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
