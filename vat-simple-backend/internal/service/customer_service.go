package service

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/models"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/repository"
	"github.com/hoanggtan02/my-app/vat-simple-backend/pkg/errors"
)

type CustomerService interface {
	CreateCustomer(req *models.CreateCustomerRequest, companyID string) (*models.Customer, *errors.AppError)
	GetCustomer(customerID, companyID string) (*models.Customer, *errors.AppError)
	ListCustomers(companyID string) ([]models.Customer, *errors.AppError)
	UpdateCustomer(customerID string, req *models.UpdateCustomerRequest, companyID string) (*models.Customer, *errors.AppError)
	DeleteCustomer(customerID, companyID string) *errors.AppError
}

type customerServiceImpl struct {
	customerRepo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerServiceImpl{customerRepo: repo}
}

func (s *customerServiceImpl) CreateCustomer(req *models.CreateCustomerRequest, companyID string) (*models.Customer, *errors.AppError) {
	customer := &models.Customer{
		Name:      req.Name,
		TaxCode:   req.TaxCode,
		Address:   req.Address,
		Email:     req.Email,
		Phone:     req.Phone,
		CompanyID: companyID,
	}

	if err := s.customerRepo.CreateCustomer(customer); err != nil {
		log.Printf("Error creating customer in repo: %v", err)
		return nil, errors.ErrInternalServerError
	}

	return customer, nil
}

func (s *customerServiceImpl) GetCustomer(customerID, companyID string) (*models.Customer, *errors.AppError) {
	customer, err := s.customerRepo.GetCustomerByID(customerID, companyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewAppError(http.StatusNotFound, "Customer not found", errors.WithCode("CUSTOMER_NOT_FOUND"))
		}
		log.Printf("Error getting customer from repo: %v", err)
		return nil, errors.ErrInternalServerError
	}
	return customer, nil
}

func (s *customerServiceImpl) ListCustomers(companyID string) ([]models.Customer, *errors.AppError) {
	customers, err := s.customerRepo.GetAllCustomers(companyID)
	if err != nil {
		log.Printf("Error listing customers from repo: %v", err)
		return nil, errors.ErrInternalServerError
	}
	return customers, nil
}

func (s *customerServiceImpl) UpdateCustomer(customerID string, req *models.UpdateCustomerRequest, companyID string) (*models.Customer, *errors.AppError) {
	// First, check if customer exists and belongs to the company
	_, appErr := s.GetCustomer(customerID, companyID)
	if appErr != nil {
		return nil, appErr
	}

	customer := &models.Customer{
		ID:        customerID,
		Name:      req.Name,
		TaxCode:   req.TaxCode,
		Address:   req.Address,
		Email:     req.Email,
		Phone:     req.Phone,
		CompanyID: companyID,
	}

	if err := s.customerRepo.UpdateCustomer(customer); err != nil {
		log.Printf("Error updating customer in repo: %v", err)
		return nil, errors.ErrInternalServerError
	}

	return customer, nil
}

func (s *customerServiceImpl) DeleteCustomer(customerID, companyID string) *errors.AppError {
	if err := s.customerRepo.DeleteCustomer(customerID, companyID); err != nil {
		log.Printf("Error deleting customer from repo: %v", err)
		return errors.ErrInternalServerError
	}
	return nil
}
