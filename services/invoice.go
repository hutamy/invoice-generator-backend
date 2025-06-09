package services

import (
	"bytes"
	"context"
	"html/template"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/hutamy/invoice-generator/dto"
	"github.com/hutamy/invoice-generator/models"
	"github.com/hutamy/invoice-generator/repositories"
)

type InvoiceService interface {
	CreateInvoice(invoice *models.Invoice) error
	GetInvoiceByID(id uint) (*models.Invoice, error)
	ListInvoiceByUserID(userID uint) ([]models.Invoice, error)
	UpdateInvoice(id uint, req *dto.UpdateInvoiceRequest) error
	GenerateInvoicePDF(invoiceID uint) ([]byte, error)
}

type invoiceService struct {
	invoiceRepo repositories.InvoiceRepository
	clientRepo  repositories.ClientRepository
	authRepo    repositories.AuthRepository
}

func NewInvoiceService(
	invoiceRepo repositories.InvoiceRepository,
	clientRepo repositories.ClientRepository,
	authRepo repositories.AuthRepository,
) InvoiceService {
	return &invoiceService{
		invoiceRepo: invoiceRepo,
		clientRepo:  clientRepo,
		authRepo:    authRepo,
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

func (s *invoiceService) GenerateInvoicePDF(invoiceID uint) ([]byte, error) {
	invoice, err := s.invoiceRepo.GetInvoiceByID(invoiceID)
	if err != nil {
		return nil, err
	}

	client, err := s.clientRepo.GetClientByID(invoice.ClientID, invoice.UserID)
	if err != nil {
		return nil, err
	}

	user, err := s.authRepo.GetUserByID(invoice.UserID)
	if err != nil {
		return nil, err
	}

	// Load HTML template
	tmpl, err := template.ParseFiles("templates/invoice.html")
	if err != nil {
		return nil, err
	}

	var htmlBuf bytes.Buffer
	err = tmpl.Execute(&htmlBuf, map[string]interface{}{
		"Invoice": invoice,
		"Client":  client,
		"User":    user,
	})
	if err != nil {
		return nil, err
	}

	// Setup headless browser
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var pdfBuf []byte
	htmlContent := htmlBuf.String()
	err = chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Set the HTML content directly
			// Set the HTML content using JavaScript evaluation
			return chromedp.Evaluate(`document.documentElement.innerHTML = `+strconv.Quote(htmlContent), nil).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, err
	}

	return pdfBuf, nil
}
