package dto

type InvoiceItemRequest struct {
	Description string  `json:"description" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required,min=1"`
	UnitPrice   float64 `json:"unit_price" validate:"required"`
}

type CreateInvoiceRequest struct {
	ClientID      uint                 `json:"client_id" validate:"required"`
	DueDate       string               `json:"due_date" validate:"required"`
	Items         []InvoiceItemRequest `json:"items" validate:"required,dive"`
	Notes         string               `json:"notes"`
	InvoiceNumber string               `json:"invoice_number" validate:"required"`
	Currency      string               `json:"currency" validate:"required,oneof=USD EUR IDR"`
	TaxRate       float64              `json:"tax_rate" validate:"required"`
}

type InvoiceItemUpdateRequest struct {
	ID          *uint   `json:"id,omitempty"`
	Description string  `json:"description" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required,min=1"`
	UnitPrice   float64 `json:"unit_price" validate:"required"`
}

type UpdateInvoiceRequest struct {
	ClientID      *uint                      `json:"client_id,omitempty"`
	DueDate       *string                    `json:"due_date,omitempty"`
	Notes         *string                    `json:"notes,omitempty"`
	Status        *string                    `json:"status,omitempty"`
	Currency      *string                    `json:"currency,omitempty"`
	TaxRate       *float64                   `json:"tax_rate,omitempty"`
	InvoiceNumber *string                    `json:"invoice_number,omitempty"`
	Items         []InvoiceItemUpdateRequest `json:"items,omitempty"`
}
