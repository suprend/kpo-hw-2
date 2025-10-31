package repository

import (
	"time"

	"kpo-hw-2/internal/domain"
)

// OperationRepository persists financial operations.
type OperationRepository interface {
	// Create persists a new operation. Returns ErrAlreadyExists on duplicate ID.
	Create(operation *domain.Operation) error
	// Update replaces an existing operation. Returns ErrNotFound if missing.
	// Implementations must persist new state for all fields except immutable ID.
	Update(operation *domain.Operation) error
	// Delete removes operation by ID. Returns ErrNotFound if missing.
	Delete(id domain.ID) error
	// Get fetches operation by ID. Returns ErrNotFound if missing.
	Get(id domain.ID) (*domain.Operation, error)
	// ListByAccountAndPeriod returns operations for account within [from, to] sorted by date ascending.
	ListByAccountAndPeriod(accountID domain.ID, from, to time.Time) ([]*domain.Operation, error)
}
