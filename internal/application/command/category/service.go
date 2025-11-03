package category

import (
	"context"

	"kpo-hw-2/internal/application/command"
	"kpo-hw-2/internal/application/facade"
	"kpo-hw-2/internal/domain"
)

type Decorators struct {
	Create []command.Decorator[*domain.Category]
	Update []command.Decorator[*domain.Category]
	Delete []command.Decorator[command.NoResult]
	List   []command.Decorator[[]*domain.Category]
	Get    []command.Decorator[*domain.Category]
}

type Service struct {
	facade     facade.CategoryFacade
	decorators Decorators
}

func NewService(f facade.CategoryFacade, decorators Decorators) *Service {
	return &Service{
		facade:     f,
		decorators: decorators,
	}
}

func (s *Service) Create(name string, typ domain.OperationType) command.Command[*domain.Category] {
	base := command.Func[*domain.Category]{
		ExecFn: func(_ context.Context) (*domain.Category, error) {
			return s.facade.CreateCategory(name, typ)
		},
		NameFn: func() string { return "category.create" },
	}
	return command.Wrap(base, s.decorators.Create...)
}

func (s *Service) Update(id domain.ID, name string, typ domain.OperationType) command.Command[*domain.Category] {
	base := command.Func[*domain.Category]{
		ExecFn: func(_ context.Context) (*domain.Category, error) {
			return s.facade.UpdateCategory(id, name, typ)
		},
		NameFn: func() string { return "category.update" },
	}
	return command.Wrap(base, s.decorators.Update...)
}

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

func (s *Service) List(typ domain.OperationType) command.Command[[]*domain.Category] {
	base := command.Func[[]*domain.Category]{
		ExecFn: func(_ context.Context) ([]*domain.Category, error) {
			return s.facade.ListCategories(typ)
		},
		NameFn: func() string { return "category.list" },
	}
	return command.Wrap(base, s.decorators.List...)
}

func (s *Service) Get(id domain.ID) command.Command[*domain.Category] {
	base := command.Func[*domain.Category]{
		ExecFn: func(_ context.Context) (*domain.Category, error) {
			return s.facade.GetCategory(id)
		},
		NameFn: func() string { return "category.get" },
	}
	return command.Wrap(base, s.decorators.Get...)
}
