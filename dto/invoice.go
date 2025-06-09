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

type GeneratePublicInvoiceRequest struct {
	InvoiceNumber string                     `json:"invoice_number" validate:"required"`
	IssueDate     string                     `json:"issue_date" validate:"required"`
	DueDate       string                     `json:"due_date" validate:"required"`
	Currency      string                     `json:"currency" validate:"required,oneof=USD EUR IDR"`
	Sender        SenderRequest              `json:"sender" validate:"required"`
	Recipient     SenderRecipientRequest     `json:"recipient" validate:"required"`
	Items         []InvoiceItemUpdateRequest `json:"items,omitempty"`
	TaxRate       *float64                   `json:"tax_rate,omitempty"`
	Notes         string                     `json:"notes"`
}

type SenderRequest struct {
	SenderRecipientRequest
	BankName          string `json:"bank_name" validate:"required"`
	BankAccountName   string `json:"bank_account_name" validate:"required"`
	BankAccountNumber string `json:"bank_account_number" validate:"required"`
}

type SenderRecipientRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Email   string `json:"email" validate:"email"`
	Phone   string `json:"phone"`
}
