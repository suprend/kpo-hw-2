package domain

import (
	"strings"
	"time"
)

type Operation struct {
	id            ID
	typ           OperationType
	bankAccountID ID
	categoryID    ID
	amount        int64
	date          time.Time
	description   string
}

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

func (o *Operation) ID() ID { return o.id }

func (o *Operation) Type() OperationType { return o.typ }

func (o *Operation) BankAccountID() ID { return o.bankAccountID }

func (o *Operation) CategoryID() ID { return o.categoryID }

func (o *Operation) Amount() int64 { return o.amount }

func (o *Operation) Date() time.Time { return o.date }

func (o *Operation) Description() string { return o.description }
