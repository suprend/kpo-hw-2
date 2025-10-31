package categories

import (
	"context"
	"strings"

	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/menus"
)

const (
	fieldCategoryName = "category_name"
	fieldCategoryType = "category_type"
)

// NewCreate constructs screen for creating a category.
func NewCreate() tui.Screen {
	var screen *menus.Screen

	validateName := func(value string) error {
		return menus.ValidateNonEmpty(value, "название не может быть пустым")
	}

	items := []menus.MenuItem{
		menus.NewInputItem(
			fieldCategoryName,
			"Название категории",
			"Введите название категории.",
			menus.InputConfig{
				Placeholder: "Например, Продукты",
			},
		),
		menus.NewSelectItem(
			fieldCategoryType,
			"Тип категории",
			"Выберите, относится ли категория к доходам или расходам.",
			[]menus.SelectOption{
				{Label: "Доход", Value: string(domain.OperationTypeIncome)},
				{Label: "Расход", Value: string(domain.OperationTypeExpense)},
			},
			menus.SelectConfig{
				InitialIndex: 0,
			},
		),
		menus.NewActionItem(
			"save",
			"Создать",
			"Сохранить категорию",
			func(ctx tui.ScreenContext, values menus.Values) tui.Result {
				name := strings.TrimSpace(values[fieldCategoryName])
				typValue := strings.TrimSpace(values[fieldCategoryType])
				hasError := menus.ApplyValidation(screen, fieldCategoryName, name, validateName)

				if typValue == "" {
					screen.SetFieldError(fieldCategoryType, "необходимо выбрать тип категории")
					hasError = true
				} else {
					screen.SetFieldError(fieldCategoryType, "")
				}

				if hasError {
					return tui.Result{}
				}

				typ := domain.OperationType(typValue)

				createCmd := ctx.CategoryCommands().Create(name, typ)
				if _, err := createCmd.Execute(context.Background()); err != nil {
					screen.SetFieldError(fieldCategoryName, err.Error())
					return tui.Result{}
				}

				menus.ClearFields(screen, fieldCategoryName)
				screen.SetValue(fieldCategoryType, string(domain.OperationTypeIncome))
				return tui.Result{
					Pop: true,
				}
			},
		),
		menus.NewPopItem("Назад", "Вернуться без сохранения"),
	}

	screen = menus.NewScreen(
		"Новая категория",
		"Создайте новую категорию операций.",
		items,
	)
	return screen
}
