package domain

import "strings"

// BankAccount represents a financial account in the domain.
type BankAccount struct {
	id      ID
	name    string
	balance int64
}

// NewBankAccount validates and constructs a bank account aggregate.
func NewBankAccount(id ID, name string, balance int64) (*BankAccount, error) {
	if id == "" {
		return nil, ErrInvalidBankAccount
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidBankAccount
	}

	return &BankAccount{
		id:      id,
		name:    name,
		balance: balance,
	}, nil
}

// ID returns the account identifier.
func (b *BankAccount) ID() ID { return b.id }

// Name returns the account name.
func (b *BankAccount) Name() string { return b.name }

// Balance returns the current balance.
func (b *BankAccount) Balance() int64 { return b.balance }

// ApplyOperation mutates account balance according to provided operation.
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

// RevertOperation rolls back balance change caused by provided operation.
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
