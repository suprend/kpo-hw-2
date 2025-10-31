package accountsmenu

import (
	"fmt"

	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/menus"
)

// NewList builds a simple screen where each account rendered as idle menu item.
func NewList(accounts []*domain.BankAccount) tui.Screen {
	items := make([]menus.Item, 0, len(accounts)+1)

	for _, account := range accounts {
		acc := account
		items = append(items, menus.NewActionItem(
			acc.ID().String(),
			acc.Name(),
			fmt.Sprintf("Баланс: %d", acc.Balance()),
			func(ctx tui.ScreenContext, values menus.Values) tui.Result {
				return tui.Result{Replace: NewEdit(acc)}
			},
		))
	}

	items = append(items, menus.NewPopItem("Назад", "Вернуться в главное меню"))

	return menus.NewScreen(
		"Список счетов",
		"Выберите счёт для редактирования или вернитесь назад.",
		items,
	).WithEmptyMessage("Счета ещё не добавлены.")
}
