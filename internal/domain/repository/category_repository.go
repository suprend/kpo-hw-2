package repository

import "kpo-hw-2/internal/domain"

// CategoryRepository persists category aggregates.
type CategoryRepository interface {
	// Create persists a new category. Returns ErrAlreadyExists on duplicate ID.
	Create(category *domain.Category) error
	// Update replaces an existing category. Returns ErrNotFound if missing.
	Update(category *domain.Category) error
	// Delete removes category by ID. Returns ErrNotFound if missing.
	Delete(id domain.ID) error
	// Get fetches category by ID. Returns ErrNotFound if missing.
	Get(id domain.ID) (*domain.Category, error)
	// ListAll returns all categories sorted by name.
	ListAll() ([]*domain.Category, error)
	// ListByType returns categories of provided type sorted by name.
	ListByType(typ domain.OperationType) ([]*domain.Category, error)
}
