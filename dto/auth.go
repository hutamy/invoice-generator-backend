package dto

type SignUpRequest struct {
	Name              string `json:"name" binding:"required"`
	Email             string `json:"email" binding:"required,email"`
	Password          string `json:"password" binding:"required,min=6"`
	Address           string `json:"address" binding:"required"`
	Phone             string `json:"phone" binding:"required"`
	BankName          string `json:"bank_name" binding:"required"`
	BankAccountName   string `json:"bank_account_name" binding:"required"`
	BankAccountNumber string `json:"bank_account_number" binding:"required"`
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
