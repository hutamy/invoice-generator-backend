package services

import (
	"bytes"
	"fmt"
	"time"

	"github.com/hutamy/invoice-generator/repositories"
	"github.com/jung-kurt/gofpdf"
)

type PDFService interface {
	GenerateInvoicePDF(invoiceID uint) ([]byte, error)
}

type pdfService struct {
	invoiceRepo repositories.InvoiceRepository
	clientRepo  repositories.ClientRepository
}

func NewPDFService(invoiceRepo repositories.InvoiceRepository, clientRepo repositories.ClientRepository) PDFService {
	return &pdfService{
		invoiceRepo: invoiceRepo,
		clientRepo:  clientRepo,
	}
}

func (s *pdfService) GenerateInvoicePDF(invoiceID uint) ([]byte, error) {
	inv, err := s.invoiceRepo.GetInvoiceByID(invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	client, err := s.clientRepo.GetClientByID(inv.ClientID, inv.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	pdf.Cell(40, 10, "Invoice")
	pdf.Ln(20)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Invoice ID: %d", inv.ID))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Client: %s", client.Name))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Amount: %.2f", inv.Subtotal))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Due Date: %s", inv.DueDate.Format("2006-01-02")))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Status: %s", inv.Status))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Created At: %s", inv.CreatedAt.Format(time.RFC1123)))

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
