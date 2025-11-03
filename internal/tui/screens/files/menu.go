package files

import (
	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/menus"
)

func NewMenu() tui.Screen {
	items := []menus.MenuItem{
		menus.NewActionItem(
			"import",
			"Импорт данных",
			"Загрузить данные из файла.",
			func(ctx tui.ScreenContext, _ menus.Values) tui.Result {
				return tui.Result{Push: newImportScreen(ctx)}
			},
		),
		menus.NewActionItem(
			"export",
			"Экспорт данных",
			"Сохранить текущие данные в файл.",
			func(ctx tui.ScreenContext, _ menus.Values) tui.Result {
				return tui.Result{Push: newExportScreen(ctx)}
			},
		),
		menus.NewPopItem("Назад", "Вернуться в главное меню"),
	}

	return menus.NewScreen(
		"Работа с файлами",
		"Выберите действие.",
		items,
	)
}
