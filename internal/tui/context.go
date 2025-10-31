package tui

import (
	accountcmd "kpo-hw-2/internal/application/command/account"
	categorycmd "kpo-hw-2/internal/application/command/category"
	operationcmd "kpo-hw-2/internal/application/command/operation"
	"kpo-hw-2/internal/application/facade"
)

// ScreenContext exposes dependencies available to screens.
type ScreenContext interface {
	Accounts() facade.AccountFacade
	Categories() facade.CategoryFacade
	Operations() facade.OperationFacade

	AccountCommands() *accountcmd.Service
	CategoryCommands() *categorycmd.Service
	OperationCommands() *operationcmd.Service
}
