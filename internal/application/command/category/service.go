package category

import (
	"context"

	"kpo-hw-2/internal/application/command"
	"kpo-hw-2/internal/application/facade"
	"kpo-hw-2/internal/domain"
)

// Decorators groups optional decorators for category commands.
type Decorators struct {
	Create []command.Decorator[*domain.Category]
	Update []command.Decorator[*domain.Category]
	Delete []command.Decorator[command.NoResult]
	List   []command.Decorator[[]*domain.Category]
	Get    []command.Decorator[*domain.Category]
}

// Service constructs commands backed by CategoryFacade.
type Service struct {
	facade     facade.CategoryFacade
	decorators Decorators
}

// NewService wires the facade and decorators for produced commands.
func NewService(f facade.CategoryFacade, decorators Decorators) *Service {
	return &Service{
		facade:     f,
		decorators: decorators,
	}
}

// Create builds a command for category creation.
func (s *Service) Create(name string, typ domain.OperationType) command.Command[*domain.Category] {
	base := command.Func[*domain.Category]{
		ExecFn: func(_ context.Context) (*domain.Category, error) {
			return s.facade.CreateCategory(name, typ)
		},
		NameFn: func() string { return "category.create" },
	}
	return command.Wrap(base, s.decorators.Create...)
}

// Update builds a command for category update.
func (s *Service) Update(id domain.ID, name string, typ domain.OperationType) command.Command[*domain.Category] {
	base := command.Func[*domain.Category]{
		ExecFn: func(_ context.Context) (*domain.Category, error) {
			return s.facade.UpdateCategory(id, name, typ)
		},
		NameFn: func() string { return "category.update" },
	}
	return command.Wrap(base, s.decorators.Update...)
}

// Delete builds a command for category deletion.
func (s *Service) Delete(id domain.ID) command.Command[command.NoResult] {
	base := command.Func[command.NoResult]{
		ExecFn: func(_ context.Context) (command.NoResult, error) {
			err := s.facade.DeleteCategory(id)
			return command.NoResult{}, err
		},
		NameFn: func() string { return "category.delete" },
	}
	return command.Wrap(base, s.decorators.Delete...)
}

// List builds a command that lists categories filtered by type.
func (s *Service) List(typ domain.OperationType) command.Command[[]*domain.Category] {
	base := command.Func[[]*domain.Category]{
		ExecFn: func(_ context.Context) ([]*domain.Category, error) {
			return s.facade.ListCategories(typ)
		},
		NameFn: func() string { return "category.list" },
	}
	return command.Wrap(base, s.decorators.List...)
}

// Get builds a command that fetches a category by ID.
func (s *Service) Get(id domain.ID) command.Command[*domain.Category] {
	base := command.Func[*domain.Category]{
		ExecFn: func(_ context.Context) (*domain.Category, error) {
			return s.facade.GetCategory(id)
		},
		NameFn: func() string { return "category.get" },
	}
	return command.Wrap(base, s.decorators.Get...)
}
