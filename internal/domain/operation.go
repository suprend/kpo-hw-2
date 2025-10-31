package domain

import (
	"strings"
	"time"
)

// Operation represents a financial transaction in the domain.
type Operation struct {
	id            ID
	typ           OperationType
	bankAccountID ID
	categoryID    ID
	amount        int64
	date          time.Time
	description   string
}

// NewOperation validates and constructs an operation aggregate.
func NewOperation(
	id ID,
	typ OperationType,
	accountID ID,
	categoryID ID,
	amount int64,
	date time.Time,
	description string,
) (*Operation, error) {
	if id == "" || accountID == "" || categoryID == "" {
		return nil, ErrInvalidOperation
	}

	if amount <= 0 {
		return nil, ErrInvalidOperation
	}

	switch typ {
	case OperationTypeIncome, OperationTypeExpense:
	default:
		return nil, ErrInvalidOperation
	}

	description = strings.TrimSpace(description)

	return &Operation{
		id:            id,
		typ:           typ,
		bankAccountID: accountID,
		categoryID:    categoryID,
		amount:        amount,
		date:          date,
		description:   description,
	}, nil
}

// ID returns operation identifier.
func (o *Operation) ID() ID { return o.id }

// Type returns operation type (income/expense).
func (o *Operation) Type() OperationType { return o.typ }

// BankAccountID returns related account id.
func (o *Operation) BankAccountID() ID { return o.bankAccountID }

// CategoryID returns related category id.
func (o *Operation) CategoryID() ID { return o.categoryID }

// Amount returns signed amount.
func (o *Operation) Amount() int64 { return o.amount }

// Date returns operation date.
func (o *Operation) Date() time.Time { return o.date }

// Description returns optional comment.
func (o *Operation) Description() string { return o.description }
