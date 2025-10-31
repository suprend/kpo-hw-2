package accountsmenu

import (
	"context"

	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/menus"
)

// NewMenu returns main accounts management menu.
func NewMenu() tui.Screen {
	items := []menus.MenuItem{
		menus.NewActionItem(
			"list",
			"Список счетов",
			"Просмотреть все счета (пока без действий).",
			func(ctx tui.ScreenContext, _ menus.Values) tui.Result {
				cmd := ctx.AccountCommands().List()
				accounts, err := cmd.Execute(context.Background())
				if err != nil {
					return tui.Result{}
				}
				return tui.Result{Push: NewList(accounts)}
			},
		),
		menus.NewActionItem(
			"create",
			"Добавить счёт",
			"Открыть форму создания нового счёта.",
			func(tui.ScreenContext, menus.Values) tui.Result {
				return tui.Result{Push: NewCreate()}
			},
		),
		menus.NewPopItem("Назад", "Вернуться в главное меню"),
	}

	return menus.NewScreen(
		"Меню счетов",
		"Выберите действие.",
		items,
	)
}
