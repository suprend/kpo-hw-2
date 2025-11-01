package bootstrap

import (
	"fmt"

	"kpo-hw-2/internal/domain"
	domainfactory "kpo-hw-2/internal/domain/factory"
	"kpo-hw-2/internal/infrastructure/di"
)

func registerDomain(container di.Container) error {
	if err := di.Register(container, func(c di.Container) (domainfactory.BankAccountFactory, error) {
		idGenerator, err := di.Resolve[domain.IDGenerator](c)
		if err != nil {
			return nil, err
		}
		return domainfactory.NewBankAccountFactory(idGenerator), nil
	}); err != nil {
		return fmt.Errorf("bootstrap: register bank account factory: %w", err)
	}

	if err := di.Register(container, func(c di.Container) (domainfactory.CategoryFactory, error) {
		idGenerator, err := di.Resolve[domain.IDGenerator](c)
		if err != nil {
			return nil, err
		}
		return domainfactory.NewCategoryFactory(idGenerator), nil
	}); err != nil {
		return fmt.Errorf("bootstrap: register category factory: %w", err)
	}

	if err := di.Register(container, func(c di.Container) (domainfactory.OperationFactory, error) {
		idGenerator, err := di.Resolve[domain.IDGenerator](c)
		if err != nil {
			return nil, err
		}
		return domainfactory.NewOperationFactory(idGenerator), nil
	}); err != nil {
		return fmt.Errorf("bootstrap: register operation factory: %w", err)
	}

	return nil
}
