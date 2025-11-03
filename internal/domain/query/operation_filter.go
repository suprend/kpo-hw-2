package query

import (
	"time"

	"kpo-hw-2/internal/domain"
)

type OperationFilter struct {
	accountID  domain.ID
	categoryID domain.ID
	typ        domain.OperationType
	from       *time.Time
	to         *time.Time
}

func NewOperationFilter() OperationFilter {
	return OperationFilter{}
}

func (f OperationFilter) ForAccount(id domain.ID) OperationFilter {
	f.accountID = id
	return f
}

func (f OperationFilter) ForCategory(id domain.ID) OperationFilter {
	f.categoryID = id
	return f
}

func (f OperationFilter) OfType(typ domain.OperationType) OperationFilter {
	f.typ = typ
	return f
}

func (f OperationFilter) From(from time.Time) OperationFilter {
	f.from = &from
	return f
}

func (f OperationFilter) To(to time.Time) OperationFilter {
	f.to = &to
	return f
}

func (f OperationFilter) Between(from, to time.Time) OperationFilter {
	f.from = &from
	f.to = &to
	return f
}

func (f OperationFilter) AccountID() domain.ID { return f.accountID }

func (f OperationFilter) CategoryID() domain.ID { return f.categoryID }

func (f OperationFilter) Type() domain.OperationType { return f.typ }

func (f OperationFilter) Period() (*time.Time, *time.Time) { return f.from, f.to }
