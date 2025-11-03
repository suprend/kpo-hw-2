package account

import (
	"context"

	"kpo-hw-2/internal/application/command"
	"kpo-hw-2/internal/application/facade"
	"kpo-hw-2/internal/domain"
)

type Decorators struct {
	Create []command.Decorator[*domain.BankAccount]
	Update []command.Decorator[*domain.BankAccount]
	Delete []command.Decorator[command.NoResult]
	List   []command.Decorator[[]*domain.BankAccount]
	Get    []command.Decorator[*domain.BankAccount]
}

type Service struct {
	facade     facade.AccountFacade
	decorators Decorators
}

func NewService(f facade.AccountFacade, decorators Decorators) *Service {
	return &Service{
		facade:     f,
		decorators: decorators,
	}
}

func (s *Service) Create(name string) command.Command[*domain.BankAccount] {
	base := command.Func[*domain.BankAccount]{
		ExecFn: func(_ context.Context) (*domain.BankAccount, error) {
			return s.facade.CreateAccount(name)
		},
		NameFn: func() string { return "account.create" },
	}
	return command.Wrap(base, s.decorators.Create...)
}

func (s *Service) Update(
	id domain.ID,
	name string,
	balance int64,
) command.Command[*domain.BankAccount] {
	base := command.Func[*domain.BankAccount]{
		ExecFn: func(_ context.Context) (*domain.BankAccount, error) {
			return s.facade.UpdateAccount(id, name, balance)
		},
		NameFn: func() string { return "account.update" },
	}
	return command.Wrap(base, s.decorators.Update...)
}

func (s *Service) Delete(id domain.ID) command.Command[command.NoResult] {
	base := command.Func[command.NoResult]{
		ExecFn: func(_ context.Context) (command.NoResult, error) {
			err := s.facade.DeleteAccount(id)
			return command.NoResult{}, err
		},
		NameFn: func() string { return "account.delete" },
	}
	return command.Wrap(base, s.decorators.Delete...)
}

func (s *Service) List() command.Command[[]*domain.BankAccount] {
	base := command.Func[[]*domain.BankAccount]{
		ExecFn: func(_ context.Context) ([]*domain.BankAccount, error) {
			return s.facade.ListAccounts()
		},
		NameFn: func() string { return "account.list" },
	}
	return command.Wrap(base, s.decorators.List...)
}

func (s *Service) Get(id domain.ID) command.Command[*domain.BankAccount] {
	base := command.Func[*domain.BankAccount]{
		ExecFn: func(_ context.Context) (*domain.BankAccount, error) {
			return s.facade.GetAccount(id)
		},
		NameFn: func() string { return "account.get" },
	}
	return command.Wrap(base, s.decorators.Get...)
}
