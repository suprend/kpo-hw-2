package repository

import "kpo-hw-2/internal/domain"

// AccountRepository persists bank account aggregates.
type AccountRepository interface {
	// Create persists a new account. Returns ErrAlreadyExists if ID already stored.
	Create(account *domain.BankAccount) error
	// Update replaces an existing account. Returns ErrNotFound if entity missing.
	Update(account *domain.BankAccount) error
	// Delete removes account by ID. Returns ErrNotFound if entity missing.
	Delete(id domain.ID) error
	// Get returns account by ID. Returns ErrNotFound if entity missing.
	Get(id domain.ID) (*domain.BankAccount, error)
	// List returns all accounts sorted by name.
	List() ([]*domain.BankAccount, error)
}
