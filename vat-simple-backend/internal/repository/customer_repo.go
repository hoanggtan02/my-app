package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/models"
)

type CustomerRepository interface {
	CreateCustomer(customer *models.Customer) error
	GetCustomerByID(customerID, companyID string) (*models.Customer, error)
	GetAllCustomers(companyID string) ([]models.Customer, error)
	FindCustomerByName(name, companyID string) (*models.Customer, error)
	UpdateCustomer(customerID, companyID string, req *models.UpdateCustomerRequest) error
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
	query := "SELECT id, name, tax_code, address, email, phone, company_id FROM customers WHERE id = ? AND company_id = ?"
	err := r.db.QueryRow(query, customerID, companyID).Scan(
		&customer.ID, &customer.Name, &customer.TaxCode, &customer.Address, &customer.Email, &customer.Phone, &customer.CompanyID,
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

	customers := make([]models.Customer, 0)
	for rows.Next() {
		var customer models.Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.TaxCode, &customer.Address, &customer.Email, &customer.Phone); err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func (r *customerRepositoryImpl) FindCustomerByName(name, companyID string) (*models.Customer, error) {
	customer := &models.Customer{}
	query := "SELECT id, name, tax_code, address, email, phone, company_id FROM customers WHERE name = ? AND company_id = ? LIMIT 1"
	err := r.db.QueryRow(query, name, companyID).Scan(
		&customer.ID, &customer.Name, &customer.TaxCode, &customer.Address, &customer.Email, &customer.Phone, &customer.CompanyID,
	)
	return customer, err
}

// --- HÀM UPDATE ĐÚNG ---
func (r *customerRepositoryImpl) UpdateCustomer(customerID, companyID string, req *models.UpdateCustomerRequest) error {
	var setParts []string
	var args []interface{}

	if req.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *req.Name)
	}
	if req.TaxCode != nil {
		setParts = append(setParts, "tax_code = ?")
		args = append(args, *req.TaxCode)
	}
	if req.Address != nil {
		setParts = append(setParts, "address = ?")
		args = append(args, *req.Address)
	}
	if req.Email != nil {
		setParts = append(setParts, "email = ?")
		args = append(args, *req.Email)
	}
	if req.Phone != nil {
		setParts = append(setParts, "phone = ?")
		args = append(args, *req.Phone)
	}

	if len(setParts) == 0 {
		return nil // Không có gì để cập nhật
	}

	query := fmt.Sprintf("UPDATE customers SET %s WHERE id = ? AND company_id = ?", strings.Join(setParts, ", "))
	args = append(args, customerID, companyID)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *customerRepositoryImpl) DeleteCustomer(customerID, companyID string) error {
	query := "DELETE FROM customers WHERE id = ? AND company_id = ?"
	_, err := r.db.Exec(query, customerID, companyID)
	return err
}
