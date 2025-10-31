package accountsmenu

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/menus"
)

const (
	fieldEditName    = "edit_name"
	fieldEditBalance = "edit_balance"
)

// NewEdit constructs screen for editing existing bank account.
func NewEdit(account *domain.BankAccount) tui.Screen {
	var screen *menus.Screen

	validateName := func(value string) error {
		return menus.ValidateNonEmpty(value, "название не может быть пустым")
	}

	validateBalance := func(value string) error {
		value = strings.TrimSpace(value)
		if value == "" {
			return errors.New("баланс не может быть пустым")
		}

		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			return errors.New("баланс должен быть числом")
		}
		return nil
	}

	items := []menus.Item{
		menus.NewInputItem(
			fieldEditName,
			"Название счёта",
			"",
			menus.InputConfig{
				Initial:  account.Name(),
				Validate: validateName,
			},
		),
		menus.NewInputItem(
			fieldEditBalance,
			"Баланс",
			"Введите целое число.",
			menus.InputConfig{
				Initial:  strconv.FormatInt(account.Balance(), 10),
				Validate: validateBalance,
			},
		),
		menus.NewActionItem(
			"save",
			"Сохранить",
			"Применить изменения.",
			func(ctx tui.ScreenContext, values menus.Values) tui.Result {
				name := strings.TrimSpace(values[fieldEditName])
				balanceStr := strings.TrimSpace(values[fieldEditBalance])
				hasError := false

				if menus.ApplyValidation(screen, fieldEditName, name, validateName) {
					hasError = true
				}

				if menus.ApplyValidation(screen, fieldEditBalance, balanceStr, validateBalance) {
					hasError = true
				}

				if hasError {
					return tui.Result{}
				}

				balance, _ := strconv.ParseInt(balanceStr, 10, 64)

				if _, err := ctx.Accounts().UpdateAccount(account.ID(), name, balance); err != nil {
					screen.SetFieldError(fieldEditName, err.Error())
					return tui.Result{}
				}

				return tui.Result{Pop: true}
			},
		),
		menus.NewActionItem(
			"delete",
			"Удалить счёт",
			"Удаление безвозвратно.",
			func(ctx tui.ScreenContext, values menus.Values) tui.Result {
				if err := ctx.Accounts().DeleteAccount(account.ID()); err != nil {
					screen.SetFieldError(fieldEditName, err.Error())
					return tui.Result{}
				}

				return tui.Result{Pop: true}
			},
		),
		menus.NewPopItem("Назад", "Вернуться к списку"),
	}

	screen = menus.NewScreen(
		fmt.Sprintf("Счёт: %s", account.Name()),
		"Обновите данные или удалите счёт.",
		items,
	)

	return screen
}
