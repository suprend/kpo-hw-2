package factory

import "kpo-hw-2/internal/domain"

type BankAccountFactory interface {
	Create(name string, balance int64) (*domain.BankAccount, error)
	Rebuild(id domain.ID, name string, balance int64) (*domain.BankAccount, error)
}

func NewBankAccountFactory(idGenerator domain.IDGenerator) BankAccountFactory {
	return &bankAccountFactory{idGenerator: idGenerator}
}

type bankAccountFactory struct {
	idGenerator domain.IDGenerator
}

func (f *bankAccountFactory) Create(name string, balance int64) (*domain.BankAccount, error) {
	id, err := f.idGenerator.NewID()
	if err != nil {
		return nil, err
	}

	return domain.NewBankAccount(id, name, balance)
}

func (f *bankAccountFactory) Rebuild(id domain.ID, name string, balance int64) (*domain.BankAccount, error) {
	return domain.NewBankAccount(id, name, balance)
}
