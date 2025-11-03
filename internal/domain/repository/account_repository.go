package repository

import "kpo-hw-2/internal/domain"

type AccountRepository interface {
	Create(account *domain.BankAccount) error
	Update(account *domain.BankAccount) error
	Delete(id domain.ID) error
	Get(id domain.ID) (*domain.BankAccount, error)
	List() ([]*domain.BankAccount, error)
}
