package facade

import (
	"time"

	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/domain/query"
)

// OperationFacade coordinates application use-cases for operations.
type OperationFacade interface {
	CreateOperation(
		typ domain.OperationType,
		accountID domain.ID,
		categoryID domain.ID,
		amount int64,
		date time.Time,
		description string,
	) (*domain.Operation, error)
	UpdateOperation(
		id domain.ID,
		typ domain.OperationType,
		accountID domain.ID,
		categoryID domain.ID,
		amount int64,
		date time.Time,
		description string,
	) (*domain.Operation, error)
	DeleteOperation(id domain.ID) error
	ListOperationsWithFilter(filter query.OperationFilter) ([]*domain.Operation, error)
	GetOperation(id domain.ID) (*domain.Operation, error)
}
