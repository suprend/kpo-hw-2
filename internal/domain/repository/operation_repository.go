package repository

import (
	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/domain/query"
)

type OperationRepository interface {
	Create(operation *domain.Operation) error
	Update(operation *domain.Operation) error
	Delete(id domain.ID) error
	Get(id domain.ID) (*domain.Operation, error)
	ListByFilter(filter query.OperationFilter) ([]*domain.Operation, error)
}
