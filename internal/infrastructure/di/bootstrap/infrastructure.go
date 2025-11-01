package bootstrap

import (
	"fmt"

	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/domain/repository"
	"kpo-hw-2/internal/infrastructure/di"
	"kpo-hw-2/internal/infrastructure/id"
	memoryrepo "kpo-hw-2/internal/infrastructure/repository/memory"
)

func registerInfrastructure(container di.Container) error {
	if err := di.Provide[domain.IDGenerator](container, id.NewULIDGenerator()); err != nil {
		return fmt.Errorf("bootstrap: provide id generator: %w", err)
	}

	if err := di.Register[repository.AccountRepository](container, func(di.Container) (repository.AccountRepository, error) {
		return memoryrepo.NewAccountRepository(), nil
	}); err != nil {
		return fmt.Errorf("bootstrap: register account repository: %w", err)
	}

	if err := di.Register[repository.CategoryRepository](container, func(di.Container) (repository.CategoryRepository, error) {
		return memoryrepo.NewCategoryRepository(), nil
	}); err != nil {
		return fmt.Errorf("bootstrap: register category repository: %w", err)
	}

	if err := di.Register[repository.OperationRepository](container, func(di.Container) (repository.OperationRepository, error) {
		return memoryrepo.NewOperationRepository(), nil
	}); err != nil {
		return fmt.Errorf("bootstrap: register operation repository: %w", err)
	}

	return nil
}
