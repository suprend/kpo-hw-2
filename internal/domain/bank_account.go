package domain

import "strings"

type BankAccount struct {
	id      ID
	name    string
	balance int64
}

func NewBankAccount(id ID, name string, balance int64) (*BankAccount, error) {
	if id == "" {
		return nil, ErrInvalidBankAccount
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidBankAccount
	}

	if balance < 0 {
		return nil, ErrInvalidBankAccount
	}

	return &BankAccount{
		id:      id,
		name:    name,
		balance: balance,
	}, nil
}

func (b *BankAccount) ID() ID { return b.id }

func (b *BankAccount) Name() string { return b.name }

func (b *BankAccount) Balance() int64 { return b.balance }

func (b *BankAccount) ApplyOperation(operation *Operation) error {
	if operation == nil {
		return ErrInvalidOperation
	}

	switch operation.Type() {
	case OperationTypeIncome:
		b.balance += operation.Amount()
		return nil
	case OperationTypeExpense:
		amount := operation.Amount()
		if amount > b.balance {
			return ErrInsufficientFunds
		}
		b.balance -= amount
		return nil
	default:
		return ErrInvalidOperation
	}
}

func (b *BankAccount) RevertOperation(operation *Operation) error {
	if operation == nil {
		return ErrInvalidOperation
	}

	switch operation.Type() {
	case OperationTypeIncome:
		amount := operation.Amount()
		if amount > b.balance {
			return ErrInvalidOperation
		}
		b.balance -= amount
		return nil
	case OperationTypeExpense:
		b.balance += operation.Amount()
		return nil
	default:
		return ErrInvalidOperation
	}
}
