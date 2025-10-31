package operations

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/menus"
)

// NewEdit constructs screen for editing an existing operation.
func NewEdit(operation *domain.Operation, accounts []*domain.BankAccount, categories []*domain.Category) tui.Screen {
	var screen *menus.Screen

	accountSelect := buildAccountSelectData(accounts)
	categorySelect := buildCategorySelectData(categories)

	initialAccountIndex := accountSelect.indexByID[operation.BankAccountID().String()]
	initialCategoryIndex := categorySelect.indexByID[operation.CategoryID().String()]
	validateName := func(value string) error {
		return menus.ValidateNonEmpty(value, "название операции не может быть пустым")
	}

	items := []menus.MenuItem{
		menus.NewInputItem(
			fieldOperationName,
			"Название операции",
			"Введите описание или назначение.",
			menus.InputConfig{
				Initial: operation.Description(),
			},
		),
		menus.NewInputItem(
			fieldOperationDate,
			"Дата",
			"Используйте формат ГГГГ-ММ-ДД.",
			menus.InputConfig{
				Initial: operation.Date().Format(dateLayout),
			},
		),
		menus.NewInputItem(
			fieldOperationAmount,
			"Сумма",
			"Укажите положительное число без разделителей.",
			menus.InputConfig{
				Initial: strconv.FormatInt(operation.Amount(), 10),
			},
		),
		menus.NewSelectItem(
			fieldOperationAccount,
			"Счёт",
			"Выберите счёт, к которому относится операция.",
			accountSelect.options,
			menus.SelectConfig{
				InitialIndex: initialAccountIndex,
			},
		),
		menus.NewSelectItem(
			fieldOperationCategory,
			"Категория",
			"Выберите категорию операции.",
			categorySelect.options,
			menus.SelectConfig{
				InitialIndex: initialCategoryIndex,
			},
		),
		menus.NewActionItem(
			"save",
			"Сохранить",
			"Применить изменения.",
			func(ctx tui.ScreenContext, values menus.Values) tui.Result {
				name := strings.TrimSpace(values[fieldOperationName])
				dateStr := strings.TrimSpace(values[fieldOperationDate])
				amountStr := strings.TrimSpace(values[fieldOperationAmount])
				accountID := strings.TrimSpace(values[fieldOperationAccount])
				categoryID := strings.TrimSpace(values[fieldOperationCategory])

				hasError := menus.ApplyValidation(screen, fieldOperationName, name, validateName)

				dateValue, dateErr := time.Parse(dateLayout, dateStr)
				if dateErr != nil {
					screen.SetFieldError(fieldOperationDate, "используйте формат ГГГГ-ММ-ДД")
					hasError = true
				} else {
					screen.SetFieldError(fieldOperationDate, "")
				}

				amountValue, amountErr := strconv.ParseInt(amountStr, 10, 64)
				if amountErr != nil || amountValue <= 0 {
					screen.SetFieldError(fieldOperationAmount, "сумма должна быть положительным числом")
					hasError = true
				} else {
					screen.SetFieldError(fieldOperationAmount, "")
				}

				if accountID == "" {
					screen.SetFieldError(fieldOperationAccount, "нужно выбрать счёт")
					hasError = true
				} else if _, exists := accountSelect.indexByID[accountID]; !exists {
					screen.SetFieldError(fieldOperationAccount, "выбранный счёт недоступен")
					hasError = true
				} else {
					screen.SetFieldError(fieldOperationAccount, "")
				}

				categoryType, categoryExists := categorySelect.typeByID[categoryID]
				if categoryID == "" {
					screen.SetFieldError(fieldOperationCategory, "нужно выбрать категорию")
					hasError = true
				} else if !categoryExists {
					screen.SetFieldError(fieldOperationCategory, "выбранная категория недоступна")
					hasError = true
				} else {
					screen.SetFieldError(fieldOperationCategory, "")
				}

				opType := operation.Type()

				if opType != "" && categoryType != "" && opType != categoryType {
					screen.SetFieldError(fieldOperationCategory, "категория не соответствует типу операции")
					hasError = true
				}

				if hasError {
					return tui.Result{}
				}

				updateCmd := ctx.OperationCommands().Update(
					operation.ID(),
					opType,
					domain.ID(accountID),
					domain.ID(categoryID),
					amountValue,
					dateValue,
					name,
				)
				if _, err := updateCmd.Execute(context.Background()); err != nil {
					switch {
					case errors.Is(err, domain.ErrInsufficientFunds):
						screen.SetFieldError(fieldOperationAmount, "на счёте недостаточно средств")
					case errors.Is(err, domain.ErrOperationTypeMismatch):
						screen.SetFieldError(fieldOperationCategory, "тип категории не совпадает с типом операции")
					default:
						screen.SetFieldError(fieldOperationName, err.Error())
					}
					return tui.Result{}
				}

				return tui.Result{Pop: true}
			},
		),
		menus.NewActionItem(
			"delete",
			"Удалить операцию",
			"Удаление безвозвратно.",
			func(ctx tui.ScreenContext, values menus.Values) tui.Result {
				deleteCmd := ctx.OperationCommands().Delete(operation.ID())
				if _, err := deleteCmd.Execute(context.Background()); err != nil {
					screen.SetFieldError(fieldOperationName, err.Error())
					return tui.Result{}
				}
				return tui.Result{Pop: true}
			},
		),
		menus.NewPopItem("Назад", "Вернуться без изменений"),
	}

	screen = menus.NewScreen(
		"Редактирование операции",
		"Обновите данные или удалите операцию.",
		items,
	)

	return screen
}
