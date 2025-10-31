package query

import (
	"time"

	"kpo-hw-2/internal/domain"
)

// OperationFilter describes criteria for selecting operations.
type OperationFilter struct {
	accountID  domain.ID
	categoryID domain.ID
	typ        domain.OperationType
	from       *time.Time
	to         *time.Time
}

// NewOperationFilter creates a filter with default values.
func NewOperationFilter() OperationFilter {
	return OperationFilter{}
}

// ForAccount restricts operations to a specific account.
func (f OperationFilter) ForAccount(id domain.ID) OperationFilter {
	f.accountID = id
	return f
}

// ForCategory restricts operations to a specific category.
func (f OperationFilter) ForCategory(id domain.ID) OperationFilter {
	f.categoryID = id
	return f
}

// OfType restricts operations by type (income/expense).
func (f OperationFilter) OfType(typ domain.OperationType) OperationFilter {
	f.typ = typ
	return f
}

// From restricts operations to those occurring on or after provided date.
func (f OperationFilter) From(from time.Time) OperationFilter {
	f.from = &from
	return f
}

// To restricts operations to those occurring on or before provided date.
func (f OperationFilter) To(to time.Time) OperationFilter {
	f.to = &to
	return f
}

// Between limits operations to a period (inclusive).
func (f OperationFilter) Between(from, to time.Time) OperationFilter {
	f.from = &from
	f.to = &to
	return f
}

// AccountID returns selected account identifier (may be zero).
func (f OperationFilter) AccountID() domain.ID { return f.accountID }

// CategoryID returns selected category identifier (may be zero).
func (f OperationFilter) CategoryID() domain.ID { return f.categoryID }

// Type returns requested operation type (may be empty).
func (f OperationFilter) Type() domain.OperationType { return f.typ }

// Period returns optional start/end dates if defined.
func (f OperationFilter) Period() (*time.Time, *time.Time) { return f.from, f.to }
