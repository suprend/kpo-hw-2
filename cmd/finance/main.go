package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	appfacade "kpo-hw-2/internal/application/facade"
	domainfactory "kpo-hw-2/internal/domain/factory"
	"kpo-hw-2/internal/infrastructure/id"
	memoryrepo "kpo-hw-2/internal/infrastructure/repository/memory"
	"kpo-hw-2/internal/tui"
	mainmenu "kpo-hw-2/internal/tui/screens/main"
)

func main() {
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

	model := tui.NewProgram(
		accountFacade,
		categoryFacade,
		operationFacade,
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
