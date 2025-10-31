package mainmenu

import (
	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/menus"
	accountsmenu "kpo-hw-2/internal/tui/screens/accounts"
	categoriesmenu "kpo-hw-2/internal/tui/screens/categories"
)

// New constructs root main menu screen.
func New() tui.Screen {
	var screen *menus.Screen

	items := []menus.MenuItem{
		menus.NewActionItem("accounts", "Счета", "Перейти к управлению счетами", func(tui.ScreenContext, menus.Values) tui.Result {
			return tui.Result{Push: accountsmenu.NewMenu()}
		}),
		menus.NewActionItem("categories", "Категории", "Управление категориями операций", func(tui.ScreenContext, menus.Values) tui.Result {
			return tui.Result{Push: categoriesmenu.NewMenu()}
		}),
		menus.NewPopItem("Выход", "Завершить работу программы"),
	}

	screen = menus.NewScreen("Главное меню", "Выберите раздел для продолжения.", items)
	return screen
}
