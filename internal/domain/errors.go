package domain

import "errors"

var (
	ErrInvalidID             = errors.New("invalid id")
	ErrInvalidBankAccount    = errors.New("invalid bank account")
	ErrInvalidCategory       = errors.New("invalid category")
	ErrInvalidOperation      = errors.New("invalid operation")
	ErrInsufficientFunds     = errors.New("insufficient funds")
	ErrOperationTypeMismatch = errors.New("operation type mismatch")
	ErrNotFound              = errors.New("not found")
	ErrAlreadyExists         = errors.New("already exists")
)
