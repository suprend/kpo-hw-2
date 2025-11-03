package operation

import (
	"context"
	"time"

	"kpo-hw-2/internal/application/command"
	"kpo-hw-2/internal/application/facade"
	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/domain/query"
)

type Decorators struct {
	Create []command.Decorator[*domain.Operation]
	Update []command.Decorator[*domain.Operation]
	Delete []command.Decorator[command.NoResult]
	List   []command.Decorator[[]*domain.Operation]
	Get    []command.Decorator[*domain.Operation]
}

type Service struct {
	facade     facade.OperationFacade
	decorators Decorators
}

func NewService(f facade.OperationFacade, decorators Decorators) *Service {
	return &Service{
		facade:     f,
		decorators: decorators,
	}
}

func (s *Service) Create(
	typ domain.OperationType,
	accountID domain.ID,
	categoryID domain.ID,
	amount int64,
	date time.Time,
	description string,
) command.Command[*domain.Operation] {
	base := command.Func[*domain.Operation]{
		ExecFn: func(_ context.Context) (*domain.Operation, error) {
			return s.facade.CreateOperation(
				typ,
				accountID,
				categoryID,
				amount,
				date,
				description,
			)
		},
		NameFn: func() string { return "operation.create" },
	}
	return command.Wrap(base, s.decorators.Create...)
}

func (s *Service) Update(
	id domain.ID,
	typ domain.OperationType,
	accountID domain.ID,
	categoryID domain.ID,
	amount int64,
	date time.Time,
	description string,
) command.Command[*domain.Operation] {
	base := command.Func[*domain.Operation]{
		ExecFn: func(_ context.Context) (*domain.Operation, error) {
			return s.facade.UpdateOperation(
				id,
				typ,
				accountID,
				categoryID,
				amount,
				date,
				description,
			)
		},
		NameFn: func() string { return "operation.update" },
	}
	return command.Wrap(base, s.decorators.Update...)
}

func (s *Service) Delete(id domain.ID) command.Command[command.NoResult] {
	base := command.Func[command.NoResult]{
		ExecFn: func(_ context.Context) (command.NoResult, error) {
			err := s.facade.DeleteOperation(id)
			return command.NoResult{}, err
		},
		NameFn: func() string { return "operation.delete" },
	}
	return command.Wrap(base, s.decorators.Delete...)
}

func (s *Service) List(filter query.OperationFilter) command.Command[[]*domain.Operation] {
	base := command.Func[[]*domain.Operation]{
		ExecFn: func(_ context.Context) ([]*domain.Operation, error) {
			return s.facade.ListOperationsWithFilter(filter)
		},
		NameFn: func() string { return "operation.list" },
	}
	return command.Wrap(base, s.decorators.List...)
}

func (s *Service) Get(id domain.ID) command.Command[*domain.Operation] {
	base := command.Func[*domain.Operation]{
		ExecFn: func(_ context.Context) (*domain.Operation, error) {
			return s.facade.GetOperation(id)
		},
		NameFn: func() string { return "operation.get" },
	}
	return command.Wrap(base, s.decorators.Get...)
}
