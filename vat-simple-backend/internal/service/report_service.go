package service

import (
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/models"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/repository"
	"github.com/hoanggtan02/my-app/vat-simple-backend/pkg/errors"
)

type ReportService interface {
	GetSalesOverTime(companyID, startDate, endDate string) ([]models.SalesDataPoint, *errors.AppError)
}

type reportServiceImpl struct {
	invoiceRepo repository.InvoiceRepository
}

func NewReportService(invoiceRepo repository.InvoiceRepository) ReportService {
	return &reportServiceImpl{invoiceRepo: invoiceRepo}
}

func (s *reportServiceImpl) GetSalesOverTime(companyID, startDate, endDate string) ([]models.SalesDataPoint, *errors.AppError) {
	data, err := s.invoiceRepo.GetSalesOverTime(companyID, startDate, endDate)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	return data, nil
}
