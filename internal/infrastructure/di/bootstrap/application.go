package bootstrap

import (
	"fmt"

	appfacade "kpo-hw-2/internal/application/facade"
	fileexport "kpo-hw-2/internal/application/files/export"
	fileimport "kpo-hw-2/internal/application/files/import"
	domainfactory "kpo-hw-2/internal/domain/factory"
	"kpo-hw-2/internal/domain/repository"
	"kpo-hw-2/internal/infrastructure/di"
)

func registerApplication(container di.Container) error {
	if err := di.Register(container, func(c di.Container) (appfacade.AccountFacade, error) {
		factory, err := di.Resolve[domainfactory.BankAccountFactory](c)
		if err != nil {
			return nil, err
		}
		repo, err := di.Resolve[repository.AccountRepository](c)
		if err != nil {
			return nil, err
		}
		return appfacade.NewAccountFacade(factory, repo), nil
	}); err != nil {
		return fmt.Errorf("bootstrap: register account facade: %w", err)
	}

	if err := di.Register(container, func(c di.Container) (appfacade.CategoryFacade, error) {
		factory, err := di.Resolve[domainfactory.CategoryFactory](c)
		if err != nil {
			return nil, err
		}
		repo, err := di.Resolve[repository.CategoryRepository](c)
		if err != nil {
			return nil, err
		}
		return appfacade.NewCategoryFacade(factory, repo), nil
	}); err != nil {
		return fmt.Errorf("bootstrap: register category facade: %w", err)
	}

	if err := di.Register(container, func(c di.Container) (appfacade.OperationFacade, error) {
		factory, err := di.Resolve[domainfactory.OperationFactory](c)
		if err != nil {
			return nil, err
		}
		opRepo, err := di.Resolve[repository.OperationRepository](c)
		if err != nil {
			return nil, err
		}
		accountRepo, err := di.Resolve[repository.AccountRepository](c)
		if err != nil {
			return nil, err
		}
		categoryRepo, err := di.Resolve[repository.CategoryRepository](c)
		if err != nil {
			return nil, err
		}
		return appfacade.NewOperationFacade(factory, opRepo, accountRepo, categoryRepo), nil
	}); err != nil {
		return fmt.Errorf("bootstrap: register operation facade: %w", err)
	}

	if err := di.Register(container, func(c di.Container) (*fileexport.Service, error) {
		accountRepo, err := di.Resolve[repository.AccountRepository](c)
		if err != nil {
			return nil, err
		}
		categoryRepo, err := di.Resolve[repository.CategoryRepository](c)
		if err != nil {
			return nil, err
		}
		operationRepo, err := di.Resolve[repository.OperationRepository](c)
		if err != nil {
			return nil, err
		}
		exporters, err := di.Resolve[[]fileexport.Exporter](c)
		if err != nil {
			return nil, err
		}

		return fileexport.NewService(accountRepo, categoryRepo, operationRepo, exporters), nil
	}); err != nil {
		return fmt.Errorf("bootstrap: register export service: %w", err)
	}

	if err := di.Register(container, func(c di.Container) (*fileimport.Service, error) {
		accountFacade, err := di.Resolve[appfacade.AccountFacade](c)
		if err != nil {
			return nil, err
		}
		categoryFacade, err := di.Resolve[appfacade.CategoryFacade](c)
		if err != nil {
			return nil, err
		}
		operationFacade, err := di.Resolve[appfacade.OperationFacade](c)
		if err != nil {
			return nil, err
		}
		importers, err := di.Resolve[[]fileimport.Importer](c)
		if err != nil {
			return nil, err
		}

		return fileimport.NewService(accountFacade, categoryFacade, operationFacade, importers), nil
	}); err != nil {
		return fmt.Errorf("bootstrap: register import service: %w", err)
	}

	return nil
}
