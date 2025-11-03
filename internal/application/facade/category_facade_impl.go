package facade

import (
	"kpo-hw-2/internal/domain"
	domainfactory "kpo-hw-2/internal/domain/factory"
	"kpo-hw-2/internal/domain/repository"
)

type categoryFacade struct {
	factory    domainfactory.CategoryFactory
	categories repository.CategoryRepository
}

func NewCategoryFacade(categoryFactory domainfactory.CategoryFactory, repo repository.CategoryRepository) CategoryFacade {
	return &categoryFacade{
		factory:    categoryFactory,
		categories: repo,
	}
}

func (f *categoryFacade) CreateCategory(name string, typ domain.OperationType) (*domain.Category, error) {
	category, err := f.factory.Create(name, typ)
	if err != nil {
		return nil, err
	}

	if err := f.categories.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (f *categoryFacade) CreateCategoryWithID(id domain.ID, name string, typ domain.OperationType) (*domain.Category, error) {
	category, err := f.factory.Rebuild(id, name, typ)
	if err != nil {
		return nil, err
	}

	if err := f.categories.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (f *categoryFacade) UpdateCategory(id domain.ID, name string, typ domain.OperationType) (*domain.Category, error) {
	category, err := f.factory.Rebuild(id, name, typ)
	if err != nil {
		return nil, err
	}

	if err := f.categories.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (f *categoryFacade) DeleteCategory(id domain.ID) error {
	if id == "" {
		return domain.ErrInvalidCategory
	}

	return f.categories.Delete(id)
}

func (f *categoryFacade) ListCategories(typ domain.OperationType) ([]*domain.Category, error) {
	switch typ {
	case "":
		return f.categories.ListAll()
	case domain.OperationTypeIncome, domain.OperationTypeExpense:
		return f.categories.ListByType(typ)
	default:
		return nil, domain.ErrInvalidCategory
	}
}

func (f *categoryFacade) GetCategory(id domain.ID) (*domain.Category, error) {
	if id == "" {
		return nil, domain.ErrInvalidCategory
	}

	return f.categories.Get(id)
}

var _ CategoryFacade = (*categoryFacade)(nil)
