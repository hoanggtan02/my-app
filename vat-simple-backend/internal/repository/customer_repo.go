package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/models"
)

type CustomerRepository interface {
	CreateCustomer(customer *models.Customer) error
	GetCustomerByID(customerID, companyID string) (*models.Customer, error)
	GetAllCustomers(companyID string) ([]models.Customer, error)
	UpdateCustomer(customer *models.Customer) error
	DeleteCustomer(customerID, companyID string) error
}

type customerRepositoryImpl struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepositoryImpl{db: db}
}

func (r *customerRepositoryImpl) CreateCustomer(customer *models.Customer) error {
	customer.ID = uuid.New().String()
	query := "INSERT INTO customers (id, name, tax_code, address, email, phone, company_id) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := r.db.Exec(query, customer.ID, customer.Name, customer.TaxCode, customer.Address, customer.Email, customer.Phone, customer.CompanyID)
	return err
}

func (r *customerRepositoryImpl) GetCustomerByID(customerID, companyID string) (*models.Customer, error) {
	customer := &models.Customer{}
	query := "SELECT id, name, tax_code, address, email, phone, company_id, created_at, updated_at FROM customers WHERE id = ? AND company_id = ?"
	err := r.db.QueryRow(query, customerID, companyID).Scan(
		&customer.ID, &customer.Name, &customer.TaxCode, &customer.Address, &customer.Email, &customer.Phone, &customer.CompanyID, &customer.CreatedAt, &customer.UpdatedAt,
	)
	return customer, err
}

func (r *customerRepositoryImpl) GetAllCustomers(companyID string) ([]models.Customer, error) {
	query := "SELECT id, name, tax_code, address, email, phone FROM customers WHERE company_id = ? ORDER BY created_at DESC"
	rows, err := r.db.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.TaxCode, &customer.Address, &customer.Email, &customer.Phone); err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func (r *customerRepositoryImpl) UpdateCustomer(customer *models.Customer) error {
	query := "UPDATE customers SET name = ?, tax_code = ?, address = ?, email = ?, phone = ? WHERE id = ? AND company_id = ?"
	_, err := r.db.Exec(query, customer.Name, customer.TaxCode, customer.Address, customer.Email, customer.Phone, customer.ID, customer.CompanyID)
	return err
}

func (r *customerRepositoryImpl) DeleteCustomer(customerID, companyID string) error {
	query := "DELETE FROM customers WHERE id = ? AND company_id = ?"
	_, err := r.db.Exec(query, customerID, companyID)
	return err
}
