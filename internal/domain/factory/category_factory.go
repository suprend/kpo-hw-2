package factory

import "kpo-hw-2/internal/domain"

type CategoryFactory interface {
	Create(name string, typ domain.OperationType) (*domain.Category, error)
	Rebuild(id domain.ID, name string, typ domain.OperationType) (*domain.Category, error)
}

func NewCategoryFactory(idGenerator domain.IDGenerator) CategoryFactory {
	return &categoryFactory{idGenerator: idGenerator}
}

type categoryFactory struct {
	idGenerator domain.IDGenerator
}

func (f *categoryFactory) Create(name string, typ domain.OperationType) (*domain.Category, error) {
	id, err := f.idGenerator.NewID()
	if err != nil {
		return nil, err
	}

	return domain.NewCategory(id, typ, name)
}

func (f *categoryFactory) Rebuild(id domain.ID, name string, typ domain.OperationType) (*domain.Category, error) {
	return domain.NewCategory(id, typ, name)
}
