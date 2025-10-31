package tui

import (
	"kpo-hw-2/internal/application/facade"
)

// ScreenContext exposes dependencies available to screens.
type ScreenContext interface {
	Accounts() facade.AccountFacade
	Categories() facade.CategoryFacade
	Operations() facade.OperationFacade
}
