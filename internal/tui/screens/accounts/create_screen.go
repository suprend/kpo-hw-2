package accountsmenu

import (
	"strings"

	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/menus"
)

const (
	fieldAccountName = "account_name"
)

func NewCreate() tui.Screen {
	var screen *menus.Screen

	validateName := func(value string) error {
		return menus.ValidateNonEmpty(value, "название не может быть пустым")
	}

	items := []menus.MenuItem{
		menus.NewInputItem(
			fieldAccountName,
			"Название счёта",
			"Введите название нового счёта.",
			menus.InputConfig{
				Placeholder: "Например, Основной",
			},
		),
		menus.NewActionItem(
			"save",
			"Создать",
			"Сохранить счёт",
			func(ctx tui.ScreenContext, values menus.Values) tui.Result {
				name := strings.TrimSpace(values[fieldAccountName])
				hasError := menus.ApplyValidation(screen, fieldAccountName, name, validateName)

				if hasError {
					return tui.Result{}
				}

				createCmd := ctx.AccountCommands().Create(name)
				if _, err := createCmd.Execute(ctx.Context()); err != nil {
					screen.SetFieldError(fieldAccountName, err.Error())
					return tui.Result{}
				}

				menus.ClearFields(screen, fieldAccountName)
				return tui.Result{
					Pop: true,
				}
			},
		),
		menus.NewPopItem("Назад", "Вернуться без сохранения"),
	}

	screen = menus.NewScreen(
		"Новый счёт",
		"Создайте новый счёт.",
		items,
	)
	return screen
}
