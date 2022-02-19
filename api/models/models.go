package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/go-ozzo/ozzo-validation/v3/is"
)

// SignUpModel ...
type SignUpModel struct {
	Username string `json:"username" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email"`
	Password string `json:"password" validate:"required"`
}

// Validate Register Model
func (rm *SignUpModel) Validate() error {
	return validation.ValidateStruct(
		rm,
		validation.Field(&rm.Email, validation.Required, is.Email),
		validation.Field(&rm.Password, validation.Required, validation.Length(8, 30), validation.Match(regexp.MustCompile("[a-z]|[A-Z][0-9]"))),
		validation.Field(&rm.Username, validation.Required, validation.Length(5, 30)),
	)
}

// SignUpResponseModel ...
type SignUpResponseModel struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// StandardErrorModel for all type of errors
type StandardErrorModel struct {
	ErrorMessage string `json:"error_message"`
}

// SignUpModelForUnidentifiedUser ...
type SignUpModelForUnidentifiedUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate Register Model
func (rm *SignUpModelForUnidentifiedUser) Validate() error {
	return validation.ValidateStruct(
		rm,
		validation.Field(&rm.Password, validation.Required, validation.Length(8, 30), validation.Match(regexp.MustCompile("[a-z]|[A-Z][0-9]"))),
		validation.Field(&rm.Username, validation.Required, validation.Length(5, 30)),
	)
}

// SignUpResponseModelForUnidentifiedUser ...
type SignUpResponseModelForUnidentifiedUser struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// CheckUserAccountResponseModel ...
type CheckUserAccountResponseModel struct {
	Exists bool `json:"exists"`
}

// IncomeModel ...
type IncomeModel struct {
	IncomeAmount int64 `json:"income_amount"`
}

// Success ...
type Success struct {
	Success bool `json:"success"`
}

// ExpenseModel ...
type ExpenseModel struct {
	ExpenseAmount int64 `json:"expense_amount"`
}

// GetBalanceResponseModel ...
type GetBalanceResponseModel struct {
	Balance int64 `json:"balance"`
}

// Operation ...
type Operation struct {
	Action string `json:"action"`
	Data   string `json:"date"`
}

// ListOperationsByTypeResponseModel ...
type ListOperationsByTypeResponseModel struct {
	Results []Operation `json:"results"`
	Count   int64       `json:"count"`
}
