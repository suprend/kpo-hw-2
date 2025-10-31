package operation

import (
	"context"
	"time"

	"kpo-hw-2/internal/application/command"
	"kpo-hw-2/internal/application/facade"
	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/domain/query"
)

// Decorators groups optional decorators for operation commands.
type Decorators struct {
	Create []command.Decorator[*domain.Operation]
	Update []command.Decorator[*domain.Operation]
	Delete []command.Decorator[command.NoResult]
	List   []command.Decorator[[]*domain.Operation]
	Get    []command.Decorator[*domain.Operation]
}

// Service constructs commands backed by OperationFacade.
type Service struct {
	facade     facade.OperationFacade
	decorators Decorators
}

// NewService wires the facade and decorators for produced commands.
func NewService(f facade.OperationFacade, decorators Decorators) *Service {
	return &Service{
		facade:     f,
		decorators: decorators,
	}
}

// Create builds a command that creates a financial operation.
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

// Update builds a command that updates an operation.
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

// Delete builds a command that removes an operation.
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

// List builds a command that lists operations using a filter.
func (s *Service) List(filter query.OperationFilter) command.Command[[]*domain.Operation] {
	base := command.Func[[]*domain.Operation]{
		ExecFn: func(_ context.Context) ([]*domain.Operation, error) {
			return s.facade.ListOperationsWithFilter(filter)
		},
		NameFn: func() string { return "operation.list" },
	}
	return command.Wrap(base, s.decorators.List...)
}

// Get builds a command that fetches an operation by ID.
func (s *Service) Get(id domain.ID) command.Command[*domain.Operation] {
	base := command.Func[*domain.Operation]{
		ExecFn: func(_ context.Context) (*domain.Operation, error) {
			return s.facade.GetOperation(id)
		},
		NameFn: func() string { return "operation.get" },
	}
	return command.Wrap(base, s.decorators.Get...)
}
