package facade

import "kpo-hw-2/internal/domain"

// AccountFacade coordinates application use-cases for bank accounts.
type AccountFacade interface {
	CreateAccount(name string) (*domain.BankAccount, error)
	UpdateAccount(id domain.ID, name string, balance int64) (*domain.BankAccount, error)
	DeleteAccount(id domain.ID) error
	ListAccounts() ([]*domain.BankAccount, error)
	GetAccount(id domain.ID) (*domain.BankAccount, error)
}
