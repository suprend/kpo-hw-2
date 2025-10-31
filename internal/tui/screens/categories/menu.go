package categories

import (
	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/menus"
)

// NewMenu constructs placeholder screen for category management.
func NewMenu() tui.Screen {
	items := []menus.MenuItem{
		menus.NewActionItem("list", "Список категорий", "Открыть список доступных категорий", func(ctx tui.ScreenContext, _ menus.Values) tui.Result {
			categories, err := ctx.Categories().ListCategories("")
			if err != nil {
				return tui.Result{}
			}
			return tui.Result{Push: NewList(categories)}
		}),
		menus.NewActionItem("create", "Создать категорию", "Добавить новую категорию", func(tui.ScreenContext, menus.Values) tui.Result {
			return tui.Result{Push: NewCreate()}
		}),
		menus.NewPopItem("Назад", "Вернуться в главное меню"),
	}

	return menus.NewScreen("Категории", "Выберите действие.", items)
}
