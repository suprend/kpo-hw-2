package factory

import (
	"time"

	"kpo-hw-2/internal/domain"
)

// OperationFactory creates and rebuilds operation aggregates.
type OperationFactory interface {
	Create(
		typ domain.OperationType,
		accountID domain.ID,
		categoryID domain.ID,
		amount int64,
		date time.Time,
		description string,
	) (*domain.Operation, error)
	Rebuild(
		id domain.ID,
		typ domain.OperationType,
		accountID domain.ID,
		categoryID domain.ID,
		amount int64,
		date time.Time,
		description string,
	) (*domain.Operation, error)
}

// NewOperationFactory constructs factory backed by ID generator.
func NewOperationFactory(idGenerator domain.IDGenerator) OperationFactory {
	return &operationFactory{idGenerator: idGenerator}
}

type operationFactory struct {
	idGenerator domain.IDGenerator
}

func (f *operationFactory) Create(
	typ domain.OperationType,
	accountID domain.ID,
	categoryID domain.ID,
	amount int64,
	date time.Time,
	description string,
) (*domain.Operation, error) {
	id, err := f.idGenerator.NewID()
	if err != nil {
		return nil, err
	}

	return f.Rebuild(id, typ, accountID, categoryID, amount, date, description)
}

func (f *operationFactory) Rebuild(
	id domain.ID,
	typ domain.OperationType,
	accountID domain.ID,
	categoryID domain.ID,
	amount int64,
	date time.Time,
	description string,
) (*domain.Operation, error) {
	return domain.NewOperation(
		id,
		typ,
		accountID,
		categoryID,
		amount,
		date,
		description,
	)
}
