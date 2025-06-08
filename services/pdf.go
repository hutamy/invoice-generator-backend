package services

import (
	"bytes"
	"context"
	"html/template"
	"strconv"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/hutamy/invoice-generator/repositories"
)

type PDFService interface {
	GenerateInvoicePDF(invoiceID uint) ([]byte, error)
}

type pdfService struct {
	invoiceRepo repositories.InvoiceRepository
	clientRepo  repositories.ClientRepository
	authRepo    repositories.AuthRepository
}

func NewPDFService(
	invoiceRepo repositories.InvoiceRepository,
	clientRepo repositories.ClientRepository,
	authRepo repositories.AuthRepository,
) PDFService {
	return &pdfService{
		invoiceRepo: invoiceRepo,
		clientRepo:  clientRepo,
		authRepo:    authRepo,
	}
}

func (s *pdfService) GenerateInvoicePDF(invoiceID uint) ([]byte, error) {
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
