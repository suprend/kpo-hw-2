package categories

import (
	"context"
	"fmt"
	"strings"

	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/menus"
)

const (
	fieldEditCategoryName = "edit_category_name"
	fieldEditCategoryType = "edit_category_type"
)

// NewEdit constructs screen for editing an existing category.
func NewEdit(category *domain.Category) tui.Screen {
	var screen *menus.Screen

	validateName := func(value string) error {
		return menus.ValidateNonEmpty(value, "название не может быть пустым")
	}

	items := []menus.MenuItem{
		menus.NewInputItem(
			fieldEditCategoryName,
			"Название категории",
			"",
			menus.InputConfig{
				Initial: category.Name(),
			},
		),
		menus.NewSelectItem(
			fieldEditCategoryType,
			"Тип категории",
			"Выберите тип категории.",
			[]menus.SelectOption{
				{Label: "Доход", Value: string(domain.OperationTypeIncome)},
				{Label: "Расход", Value: string(domain.OperationTypeExpense)},
			},
			menus.SelectConfig{
				InitialIndex: initialTypeIndex(category.Type()),
			},
		),
		menus.NewActionItem(
			"save",
			"Сохранить",
			"Применить изменения.",
			func(ctx tui.ScreenContext, values menus.Values) tui.Result {
				name := strings.TrimSpace(values[fieldEditCategoryName])
				typValue := strings.TrimSpace(values[fieldEditCategoryType])
				hasError := menus.ApplyValidation(screen, fieldEditCategoryName, name, validateName)

				if typValue == "" {
					screen.SetFieldError(fieldEditCategoryType, "нужно выбрать тип категории")
					hasError = true
				} else {
					screen.SetFieldError(fieldEditCategoryType, "")
				}

				if hasError {
					return tui.Result{}
				}

				typ := domain.OperationType(typValue)
				updateCmd := ctx.CategoryCommands().Update(category.ID(), name, typ)
				if _, err := updateCmd.Execute(context.Background()); err != nil {
					screen.SetFieldError(fieldEditCategoryName, err.Error())
					return tui.Result{}
				}

				listCmd := ctx.CategoryCommands().List("")
				categories, err := listCmd.Execute(context.Background())
				if err != nil {
					return tui.Result{}
				}
				return tui.Result{Replace: NewList(categories)}
			},
		),
		menus.NewActionItem(
			"delete",
			"Удалить категорию",
			"Удаление безвозвратно.",
			func(ctx tui.ScreenContext, values menus.Values) tui.Result {
				deleteCmd := ctx.CategoryCommands().Delete(category.ID())
				if _, err := deleteCmd.Execute(context.Background()); err != nil {
					screen.SetFieldError(fieldEditCategoryName, err.Error())
					return tui.Result{}
				}

				listCmd := ctx.CategoryCommands().List("")
				categories, err := listCmd.Execute(context.Background())
				if err != nil {
					return tui.Result{}
				}
				return tui.Result{Replace: NewList(categories)}
			},
		),
		menus.NewPopItem("Назад", "Вернуться к списку категорий"),
	}

	screen = menus.NewScreen(
		fmt.Sprintf("Категория: %s", category.Name()),
		"Обновите данные или удалите категорию.",
		items,
	)

	return screen
}

func initialTypeIndex(typ domain.OperationType) int {
	switch typ {
	case domain.OperationTypeExpense:
		return 1
	default:
		return 0
	}
}
