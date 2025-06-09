package dto

type CreateClientRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Address string `json:"address" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	UserID  uint   `json:"-"`
}

type UpdateClientRequest struct {
	Name    *string `json:"name" validate:"omitempty"`
	Email   *string `json:"email" validate:"omitempty,email"`
	Address *string `json:"address" validate:"omitempty"`
	Phone   *string `json:"phone" validate:"omitempty"`
	UserID  uint    `json:"-"`
	ID      uint    `param:"id" validate:"required"`
}
