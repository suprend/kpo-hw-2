package bootstrap

import (
	"fmt"
	"time"

	"kpo-hw-2/internal/application/command"
	accountcmd "kpo-hw-2/internal/application/command/account"
	categorycmd "kpo-hw-2/internal/application/command/category"
	"kpo-hw-2/internal/application/command/decorator"
	exportcmd "kpo-hw-2/internal/application/command/export"
	fileimportcmd "kpo-hw-2/internal/application/command/import"
	operationcmd "kpo-hw-2/internal/application/command/operation"
	appfacade "kpo-hw-2/internal/application/facade"
	appfiles "kpo-hw-2/internal/application/files"
	fileexport "kpo-hw-2/internal/application/files/export"
	fileimport "kpo-hw-2/internal/application/files/import"
	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/infrastructure/di"
)

func registerCommands(container di.Container) error {
	if err := di.Register(container, func(c di.Container) (*accountcmd.Service, error) {
		facade, err := di.Resolve[appfacade.AccountFacade](c)
		if err != nil {
			return nil, err
		}
		logFn, err := di.Resolve[func(string, time.Duration, error)](c)
		if err != nil {
			return nil, err
		}

		timedBankAccount := decorator.Timed[*domain.BankAccount]{Log: logFn}
		timedNoResult := decorator.Timed[command.NoResult]{Log: logFn}
		timedList := decorator.Timed[[]*domain.BankAccount]{Log: logFn}

		return accountcmd.NewService(
			facade,
			accountcmd.Decorators{
				Create: []command.Decorator[*domain.BankAccount]{timedBankAccount},
				Update: []command.Decorator[*domain.BankAccount]{timedBankAccount},
				Delete: []command.Decorator[command.NoResult]{timedNoResult},
				List:   []command.Decorator[[]*domain.BankAccount]{timedList},
				Get:    []command.Decorator[*domain.BankAccount]{timedBankAccount},
			},
		), nil
	}); err != nil {
		return fmt.Errorf("bootstrap: register account commands: %w", err)
	}

	if err := di.Register(container, func(c di.Container) (*categorycmd.Service, error) {
		facade, err := di.Resolve[appfacade.CategoryFacade](c)
		if err != nil {
			return nil, err
		}
		logFn, err := di.Resolve[func(string, time.Duration, error)](c)
		if err != nil {
			return nil, err
		}

		timedCategory := decorator.Timed[*domain.Category]{Log: logFn}
		timedNoResult := decorator.Timed[command.NoResult]{Log: logFn}
		timedList := decorator.Timed[[]*domain.Category]{Log: logFn}

		return categorycmd.NewService(
			facade,
			categorycmd.Decorators{
				Create: []command.Decorator[*domain.Category]{timedCategory},
				Update: []command.Decorator[*domain.Category]{timedCategory},
				Delete: []command.Decorator[command.NoResult]{timedNoResult},
				List:   []command.Decorator[[]*domain.Category]{timedList},
				Get:    []command.Decorator[*domain.Category]{timedCategory},
			},
		), nil
	}); err != nil {
		return fmt.Errorf("bootstrap: register category commands: %w", err)
	}

	if err := di.Register(container, func(c di.Container) (*operationcmd.Service, error) {
		facade, err := di.Resolve[appfacade.OperationFacade](c)
		if err != nil {
			return nil, err
		}
		logFn, err := di.Resolve[func(string, time.Duration, error)](c)
		if err != nil {
			return nil, err
		}

		timedOperation := decorator.Timed[*domain.Operation]{Log: logFn}
		timedNoResult := decorator.Timed[command.NoResult]{Log: logFn}
		timedList := decorator.Timed[[]*domain.Operation]{Log: logFn}

		return operationcmd.NewService(
			facade,
			operationcmd.Decorators{
				Create: []command.Decorator[*domain.Operation]{timedOperation},
				Update: []command.Decorator[*domain.Operation]{timedOperation},
				Delete: []command.Decorator[command.NoResult]{timedNoResult},
				List:   []command.Decorator[[]*domain.Operation]{timedList},
				Get:    []command.Decorator[*domain.Operation]{timedOperation},
			},
		), nil
	}); err != nil {
		return fmt.Errorf("bootstrap: register operation commands: %w", err)
	}

	if err := di.Register(container, func(c di.Container) (*exportcmd.Service, error) {
		service, err := di.Resolve[*fileexport.Service](c)
		if err != nil {
			return nil, err
		}
		logFn, err := di.Resolve[func(string, time.Duration, error)](c)
		if err != nil {
			return nil, err
		}

		timedFormats := decorator.Timed[[]appfiles.Format]{Log: logFn}
		timedNoResult := decorator.Timed[command.NoResult]{Log: logFn}

		return exportcmd.NewService(
			service,
			exportcmd.Decorators{
				ListFormats:  []command.Decorator[[]appfiles.Format]{timedFormats},
				Export:       []command.Decorator[command.NoResult]{timedNoResult},
				ExportToPath: []command.Decorator[command.NoResult]{timedNoResult},
			},
		), nil
	}); err != nil {
		return fmt.Errorf("bootstrap: register export commands: %w", err)
	}

	if err := di.Register(container, func(c di.Container) (*fileimportcmd.Service, error) {
		service, err := di.Resolve[*fileimport.Service](c)
		if err != nil {
			return nil, err
		}
		logFn, err := di.Resolve[func(string, time.Duration, error)](c)
		if err != nil {
			return nil, err
		}

		timedFormats := decorator.Timed[[]appfiles.Format]{Log: logFn}
		timedResult := decorator.Timed[fileimport.Result]{Log: logFn}

		return fileimportcmd.NewService(
			service,
			fileimportcmd.Decorators{
				ListFormats:    []command.Decorator[[]appfiles.Format]{timedFormats},
				ImportFromPath: []command.Decorator[fileimport.Result]{timedResult},
			},
		), nil
	}); err != nil {
		return fmt.Errorf("bootstrap: register import commands: %w", err)
	}

	return nil
}
