package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"kpo-hw-2/internal/application/command"
	accountcmd "kpo-hw-2/internal/application/command/account"
	categorycmd "kpo-hw-2/internal/application/command/category"
	"kpo-hw-2/internal/application/command/decorator"
	operationcmd "kpo-hw-2/internal/application/command/operation"
	appfacade "kpo-hw-2/internal/application/facade"
	"kpo-hw-2/internal/domain"
	domainfactory "kpo-hw-2/internal/domain/factory"
	"kpo-hw-2/internal/infrastructure/id"
	memoryrepo "kpo-hw-2/internal/infrastructure/repository/memory"
	"kpo-hw-2/internal/tui"
	mainmenu "kpo-hw-2/internal/tui/screens/main"
)

func main() {
	logFn, closeLog, err := openTimingLogger("logs/timings.log")
	if err != nil {
		log.Fatalf("не удалось открыть лог таймингов: %v", err)
	}
	defer func() {
		if err := closeLog(); err != nil {
			log.Printf("не удалось закрыть лог таймингов: %v", err)
		}
	}()

	// shared dependencies
	idGenerator := id.NewULIDGenerator()

	accountFactory := domainfactory.NewBankAccountFactory(idGenerator)
	categoryFactory := domainfactory.NewCategoryFactory(idGenerator)
	operationFactory := domainfactory.NewOperationFactory(idGenerator)

	accountRepository := memoryrepo.NewAccountRepository()
	categoryRepository := memoryrepo.NewCategoryRepository()
	operationRepository := memoryrepo.NewOperationRepository()

	accountFacade := appfacade.NewAccountFacade(accountFactory, accountRepository)
	categoryFacade := appfacade.NewCategoryFacade(categoryFactory, categoryRepository)
	operationFacade := appfacade.NewOperationFacade(
		operationFactory,
		operationRepository,
		accountRepository,
		categoryRepository,
	)

	rootScreen := mainmenu.New()
	if rootScreen == nil {
		log.Fatal("главное меню не инициализировано")
	}

	accountCommands := newAccountCommands(accountFacade, logFn)
	categoryCommands := newCategoryCommands(categoryFacade, logFn)
	operationCommands := newOperationCommands(operationFacade, logFn)

	model := tui.NewProgram(
		accountCommands,
		categoryCommands,
		operationCommands,
		rootScreen,
	)

	if err := runProgram(model); err != nil {
		log.Fatalf("не удалось запустить интерфейс: %v", err)
	}
}

func runProgram(model tea.Model) error {
	program := tea.NewProgram(model, tea.WithAltScreen())
	_, err := program.Run()
	return err
}

func openTimingLogger(path string) (func(name string, duration time.Duration, err error), func() error, error) {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, nil, err
		}
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, nil, err
	}

	logger := log.New(file, "", log.LstdFlags)
	logFn := func(name string, duration time.Duration, cmdErr error) {
		ms := float64(duration) / float64(time.Millisecond)
		logger.Printf("%s took %.3fms (err=%v)", name, ms, cmdErr)
	}

	closeFn := func() error {
		return file.Close()
	}

	return logFn, closeFn, nil
}

func newAccountCommands(
	facade appfacade.AccountFacade,
	logFn func(string, time.Duration, error),
) *accountcmd.Service {
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
	)
}

func newCategoryCommands(
	facade appfacade.CategoryFacade,
	logFn func(string, time.Duration, error),
) *categorycmd.Service {
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
	)
}

func newOperationCommands(
	facade appfacade.OperationFacade,
	logFn func(string, time.Duration, error),
) *operationcmd.Service {
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
	)
}
