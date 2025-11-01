package operations

import (
	"fmt"
	"strings"

	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/menus"
)

// NewMenu constructs placeholder screen for operations management.
func NewMenu() tui.Screen {
	items := []menus.MenuItem{
		menus.NewActionItem("list", "Список операций", "Просмотреть операции с фильтрами.", func(ctx tui.ScreenContext, _ menus.Values) tui.Result {
			accountCmd := ctx.AccountCommands().List()
			accounts, err := accountCmd.Execute(ctx.Context())
			if err != nil {
				return tui.Result{Push: errorScreen("Ошибка", fmt.Sprintf("Не удалось получить список счетов:\n%s", err.Error()))}
			}

			categoryCmd := ctx.CategoryCommands().List("")
			categories, err := categoryCmd.Execute(ctx.Context())
			if err != nil {
				return tui.Result{Push: errorScreen("Ошибка", fmt.Sprintf("Не удалось получить список категорий:\n%s", err.Error()))}
			}

			return tui.Result{Push: NewFilter(accounts, categories)}
		}),
		menus.NewActionItem("create", "Добавить операцию", "Создать новую финансовую операцию.", func(ctx tui.ScreenContext, _ menus.Values) tui.Result {
			accountCmd := ctx.AccountCommands().List()
			accounts, err := accountCmd.Execute(ctx.Context())
			if err != nil {
				return tui.Result{Push: errorScreen("Ошибка", fmt.Sprintf("Не удалось получить список счетов:\n%s", err.Error()))}
			}

			categoryCmd := ctx.CategoryCommands().List("")
			categories, err := categoryCmd.Execute(ctx.Context())
			if err != nil {
				return tui.Result{Push: errorScreen("Ошибка", fmt.Sprintf("Не удалось получить список категорий:\n%s", err.Error()))}
			}

			if len(accounts) == 0 || len(categories) == 0 {
				var b strings.Builder
				b.WriteString("Для создания операции необходимо:\n")
				if len(accounts) == 0 {
					b.WriteString("- Добавить хотя бы один счёт.\n")
				}
				if len(categories) == 0 {
					b.WriteString("- Создать хотя бы одну категорию.\n")
				}

				return tui.Result{Push: errorScreen("Недостаточно данных", b.String())}
			}

			return tui.Result{Push: NewCreate(accounts, categories)}
		}),
		menus.NewPopItem("Назад", "Вернуться в главное меню"),
	}

	return menus.NewScreen("Операции", "Выберите действие.", items)
}

func errorScreen(title, message string) tui.Screen {
	return menus.NewScreen(
		title,
		message,
		[]menus.MenuItem{
			menus.NewPopItem("Назад", "Вернуться в меню операций"),
		},
	)
}
