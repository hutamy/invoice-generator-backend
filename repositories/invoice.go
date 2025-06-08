package repositories

import (
	"fmt"
	"time"

	"github.com/hutamy/invoice-generator/dto"
	"github.com/hutamy/invoice-generator/models"
	"github.com/hutamy/invoice-generator/utils/errors"
	"gorm.io/gorm"
)

type InvoiceRepository interface {
	CreateInvoice(invoice *models.Invoice) error
	GetInvoiceByID(id uint) (*models.Invoice, error)
	ListInvoiceByUserID(userID uint) ([]models.Invoice, error)
	UpdateInvoice(id uint, req *dto.UpdateInvoiceRequest) error
}

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) CreateInvoice(invoice *models.Invoice) error {
	return r.db.Create(invoice).Error
}

func (r *invoiceRepository) GetInvoiceByID(id uint) (*models.Invoice, error) {
	var invoice models.Invoice
	if err := r.db.Preload("Items").First(&invoice, id).Error; err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (r *invoiceRepository) ListInvoiceByUserID(userID uint) ([]models.Invoice, error) {
	var invoices []models.Invoice
	if err := r.db.Where("user_id = ?", userID).Preload("Items").Find(&invoices).Error; err != nil {
		return nil, err
	}

	return invoices, nil
}

func (r *invoiceRepository) UpdateInvoice(id uint, req *dto.UpdateInvoiceRequest) error {
	var invoice models.Invoice
	if err := r.db.Preload("Items").First(&invoice, id).Error; err != nil {
		return err
	}

	// Update simple fields if present
	if req.ClientID != nil {
		invoice.ClientID = *req.ClientID
	}
	if req.DueDate != nil {
		dueDate, err := time.Parse(time.DateOnly, *req.DueDate)
		if err != nil {
			return errors.ErrInvalidDateFormat
		}

		invoice.DueDate = dueDate
	}
	if req.Notes != nil {
		invoice.Notes = *req.Notes
	}
	if req.Status != nil {
		invoice.Status = *req.Status
	}

	if req.Currency != nil {
		invoice.Currency = *req.Currency
	}

	if req.TaxRate != nil {
		invoice.TaxRate = *req.TaxRate
	}

	if req.InvoiceNumber != nil {
		invoice.InvoiceNumber = *req.InvoiceNumber
	}

	// Map existing items by ID
	existingItems := map[uint]models.InvoiceItem{}
	for _, item := range invoice.Items {
		existingItems[item.ID] = item
	}

	// Track IDs from request to keep
	var idsToKeep []uint
	var subtotal float64
	for _, itemReq := range req.Items {
		if itemReq.ID != nil {
			// Update existing item
			if existingItem, ok := existingItems[*itemReq.ID]; ok {
				existingItem.Description = itemReq.Description
				existingItem.Quantity = itemReq.Quantity
				existingItem.UnitPrice = itemReq.UnitPrice
				existingItem.Total = float64(itemReq.Quantity) * itemReq.UnitPrice
				subtotal += existingItem.Total
				if err := r.db.Save(&existingItem).Error; err != nil {
					return err
				}
				idsToKeep = append(idsToKeep, *itemReq.ID)
			} else {
				// ID not found in DB, return error or ignore
				return fmt.Errorf("invoice item with ID %d not found", *itemReq.ID)
			}
		} else {
			// New item to create
			newItem := models.InvoiceItem{
				InvoiceID:   invoice.ID,
				Description: itemReq.Description,
				Quantity:    itemReq.Quantity,
				UnitPrice:   itemReq.UnitPrice,
			}
			newItem.Total = float64(itemReq.Quantity) * itemReq.UnitPrice
			subtotal += newItem.Total
			if err := r.db.Create(&newItem).Error; err != nil {
				return err
			}
			idsToKeep = append(idsToKeep, newItem.ID)
		}
	}

	// Delete items not in idsToKeep
	for _, existingItem := range invoice.Items {
		found := false
		for _, id := range idsToKeep {
			if existingItem.ID == id {
				found = true
				break
			}
		}
		if !found {
			if err := r.db.Delete(&existingItem).Error; err != nil {
				return err
			}
		}
	}

	invoice.Subtotal = subtotal
	invoice.Tax = invoice.TaxRate * subtotal / 100
	invoice.Total = invoice.Subtotal + invoice.Tax
	return r.db.Save(&invoice).Error
}
