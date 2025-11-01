package operations

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/menus"
)

const (
	fieldOperationName     = "operation_name"
	fieldOperationDate     = "operation_date"
	fieldOperationAmount   = "operation_amount"
	fieldOperationAccount  = "operation_account"
	fieldOperationCategory = "operation_category"
	fieldOperationType     = "operation_type"
)

// NewCreate constructs screen for creating a financial operation.
func NewCreate(accounts []*domain.BankAccount, categories []*domain.Category) tui.Screen {
	var screen *menus.Screen

	accountSelect := buildAccountSelectData(accounts)
	categorySelect := buildCategorySelectData(categories)

	initialAccountIndex := 0
	if len(accountSelect.options) == 0 {
		initialAccountIndex = -1
	}

	initialCategoryIndex := 0
	if len(categorySelect.options) == 0 {
		initialCategoryIndex = -1
	}

	validateName := func(value string) error {
		return menus.ValidateNonEmpty(value, "название операции не может быть пустым")
	}

	items := []menus.MenuItem{
		menus.NewInputItem(
			fieldOperationName,
			"Название операции",
			"Введите описание или назначение.",
			menus.InputConfig{
				Placeholder: "Например, Зарплата",
			},
		),
		menus.NewInputItem(
			fieldOperationDate,
			"Дата",
			"Используйте формат ГГГГ-ММ-ДД.",
			menus.InputConfig{
				Initial: time.Now().Format(dateLayout),
			},
		),
		menus.NewInputItem(
			fieldOperationAmount,
			"Сумма",
			"Укажите положительное число без разделителей.",
			menus.InputConfig{
				Placeholder: "Например, 15000",
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
			"Создать",
			"Сохранить операцию.",
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

				if hasError {
					return tui.Result{}
				}

				createCmd := ctx.OperationCommands().Create(
					categoryType,
					domain.ID(accountID),
					domain.ID(categoryID),
					amountValue,
					dateValue,
					name,
				)
				if _, err := createCmd.Execute(ctx.Context()); err != nil {
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
		menus.NewPopItem("Назад", "Вернуться без сохранения"),
	}

	screen = menus.NewScreen(
		"Новая операция",
		"Создайте финансовую операцию.",
		items,
	)

	return screen
}

func readableType(typ domain.OperationType) string {
	switch typ {
	case domain.OperationTypeIncome:
		return "доход"
	case domain.OperationTypeExpense:
		return "расход"
	default:
		return string(typ)
	}
}
