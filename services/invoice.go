package services

import (
	"time"

	"github.com/hutamy/invoice-generator/dto"
	"github.com/hutamy/invoice-generator/models"
	"github.com/hutamy/invoice-generator/repositories"
)

type InvoiceService interface {
	CreateInvoice(invoice *models.Invoice) error
	GetInvoiceByID(id uint) (*models.Invoice, error)
	ListInvoiceByUserID(userID uint) ([]models.Invoice, error)
	UpdateInvoice(id uint, req *dto.UpdateInvoiceRequest) error
}

type invoiceService struct {
	invoiceRepo repositories.InvoiceRepository
}

func NewInvoiceService(invoiceRepo repositories.InvoiceRepository) InvoiceService {
	return &invoiceService{
		invoiceRepo: invoiceRepo,
	}
}

func (s *invoiceService) CreateInvoice(invoice *models.Invoice) error {
	var subtotal float64
	for i, item := range invoice.Items {
		total := float64(item.Quantity) * item.UnitPrice
		invoice.Items[i].Total = total
		subtotal += total
	}

	invoice.Subtotal = subtotal
	invoice.Tax = invoice.TaxRate * subtotal / 100
	invoice.Total = invoice.Subtotal + invoice.Tax
	invoice.Status = "draft" // Default status for new invoices
	invoice.IssueDate = time.Now()
	return s.invoiceRepo.CreateInvoice(invoice)
}

func (s *invoiceService) GetInvoiceByID(id uint) (*models.Invoice, error) {
	return s.invoiceRepo.GetInvoiceByID(id)
}

func (s *invoiceService) ListInvoiceByUserID(userID uint) ([]models.Invoice, error) {
	return s.invoiceRepo.ListInvoiceByUserID(userID)
}

func (s *invoiceService) UpdateInvoice(id uint, req *dto.UpdateInvoiceRequest) error {
	return s.invoiceRepo.UpdateInvoice(id, req)
}
