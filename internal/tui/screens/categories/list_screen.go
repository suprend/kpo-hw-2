package categories

import (
	"fmt"

	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/menus"
)

func NewList(categories []*domain.Category) tui.Screen {
	items := make([]menus.MenuItem, 0, len(categories)+1)

	for _, category := range categories {
		cat := category
		items = append(items, menus.NewActionItem(
			cat.ID().String(),
			cat.Name(),
			fmt.Sprintf("Тип: %s", readableType(cat.Type())),
			func(ctx tui.ScreenContext, values menus.Values) tui.Result {
				return tui.Result{Replace: NewEdit(cat)}
			},
		))
	}

	items = append(items, menus.NewPopItem("Назад", "Вернуться к меню категорий"))

	return menus.NewScreen(
		"Список категорий",
		"Просмотрите доступные категории.",
		items,
	).WithEmptyMessage("Категории ещё не добавлены.")
}

func readableType(typ domain.OperationType) string {
	switch typ {
	case domain.OperationTypeIncome:
		return "Доход"
	case domain.OperationTypeExpense:
		return "Расход"
	default:
		return string(typ)
	}
}
