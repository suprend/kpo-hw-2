package facade

import (
	"kpo-hw-2/internal/domain"
	domainfactory "kpo-hw-2/internal/domain/factory"
	"kpo-hw-2/internal/domain/repository"
)

// accountFacade is the default AccountFacade implementation.
type accountFacade struct {
	factory  domainfactory.BankAccountFactory
	accounts repository.AccountRepository
}

// NewAccountFacade wires dependencies for account use-cases.
func NewAccountFacade(
	accountFactory domainfactory.BankAccountFactory,
	accountRepo repository.AccountRepository,
) AccountFacade {
	return &accountFacade{
		factory:  accountFactory,
		accounts: accountRepo,
	}
}

func (f *accountFacade) CreateAccount(name string) (*domain.BankAccount, error) {
	account, err := f.factory.Create(name, 0)
	if err != nil {
		return nil, err
	}

	if err := f.accounts.Create(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (f *accountFacade) UpdateAccount(id domain.ID, name string, balance int64) (*domain.BankAccount, error) {
	account, err := f.factory.Rebuild(id, name, balance)
	if err != nil {
		return nil, err
	}

	if err := f.accounts.Update(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (f *accountFacade) DeleteAccount(id domain.ID) error {
	if id == "" {
		return domain.ErrInvalidBankAccount
	}

	return f.accounts.Delete(id)
}

func (f *accountFacade) ListAccounts() ([]*domain.BankAccount, error) {
	return f.accounts.List()
}

func (f *accountFacade) GetAccount(id domain.ID) (*domain.BankAccount, error) {
	if id == "" {
		return nil, domain.ErrInvalidBankAccount
	}

	return f.accounts.Get(id)
}

// compile-time check
var _ AccountFacade = (*accountFacade)(nil)
