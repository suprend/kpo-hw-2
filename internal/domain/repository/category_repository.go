package repository

import "kpo-hw-2/internal/domain"

type CategoryRepository interface {
	Create(category *domain.Category) error
	Update(category *domain.Category) error
	Delete(id domain.ID) error
	Get(id domain.ID) (*domain.Category, error)
	ListAll() ([]*domain.Category, error)
	ListByType(typ domain.OperationType) ([]*domain.Category, error)
}
