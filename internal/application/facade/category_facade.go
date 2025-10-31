package facade

import "kpo-hw-2/internal/domain"

// CategoryFacade coordinates application use-cases for categories.
type CategoryFacade interface {
	CreateCategory(name string, typ domain.OperationType) (*domain.Category, error)
	UpdateCategory(id domain.ID, name string, typ domain.OperationType) (*domain.Category, error)
	DeleteCategory(id domain.ID) error
	ListCategories(typ domain.OperationType) ([]*domain.Category, error)
	GetCategory(id domain.ID) (*domain.Category, error)
}
