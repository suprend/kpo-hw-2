package operations

import (
	"strings"
	"time"

	appanalytics "kpo-hw-2/internal/application/analytics"
	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/domain/query"
	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/menus"
)

const (
	fieldFilterStartDate = "filter_start_date"
	fieldFilterEndDate   = "filter_end_date"
	fieldFilterAccount   = "filter_account"
	fieldFilterCategory  = "filter_category"
	fieldFilterType      = "filter_type"
)

func NewFilter(accounts []*domain.BankAccount, categories []*domain.Category) tui.Screen {
	var screen *menus.Screen

	accountOptions := make([]menus.SelectOption, 0, len(accounts)+1)
	accountOptions = append(accountOptions, menus.SelectOption{
		Label: "Все счета",
		Value: "",
	})
	for _, account := range accounts {
		accountOptions = append(accountOptions, menus.SelectOption{
			Label: account.Name(),
			Value: account.ID().String(),
		})
	}

	categoryOptions := make([]menus.SelectOption, 0, len(categories)+1)
	categoryOptions = append(categoryOptions, menus.SelectOption{
		Label: "Все категории",
		Value: "",
	})
	for _, category := range categories {
		categoryOptions = append(categoryOptions, menus.SelectOption{
			Label: category.Name(),
			Value: category.ID().String(),
		})
	}

	typeOptions := []menus.SelectOption{
		{Label: "Все типы", Value: ""},
		{Label: "Доход", Value: string(domain.OperationTypeIncome)},
		{Label: "Расход", Value: string(domain.OperationTypeExpense)},
	}

	items := []menus.MenuItem{
		menus.NewInputItem(
			fieldFilterStartDate,
			"Дата начала",
			"Оставьте пустым, чтобы не ограничивать начало периода.",
			menus.InputConfig{
				Placeholder: "ГГГГ-ММ-ДД",
			},
		),
		menus.NewInputItem(
			fieldFilterEndDate,
			"Дата окончания",
			"Оставьте пустым, чтобы не ограничивать конец периода.",
			menus.InputConfig{
				Placeholder: "ГГГГ-ММ-ДД",
			},
		),
		menus.NewSelectItem(
			fieldFilterAccount,
			"Счёт",
			"Выберите счёт или оставьте «Все счета».",
			accountOptions,
			menus.SelectConfig{InitialIndex: 0},
		),
		menus.NewSelectItem(
			fieldFilterCategory,
			"Категория",
			"Выберите категорию или оставьте «Все категории».",
			categoryOptions,
			menus.SelectConfig{InitialIndex: 0},
		),
		menus.NewSelectItem(
			fieldFilterType,
			"Тип операции",
			"Выберите приход или расход, либо оставьте «Все типы».",
			typeOptions,
			menus.SelectConfig{InitialIndex: 0},
		),
		menus.NewActionItem(
			"apply",
			"Показать операции",
			"Применить фильтр и перейти к списку.",
			func(ctx tui.ScreenContext, values menus.Values) tui.Result {
				startStr := strings.TrimSpace(values[fieldFilterStartDate])
				endStr := strings.TrimSpace(values[fieldFilterEndDate])
				accountID := strings.TrimSpace(values[fieldFilterAccount])
				categoryID := strings.TrimSpace(values[fieldFilterCategory])
				operationType := strings.TrimSpace(values[fieldFilterType])

				var startDate, endDate *time.Time
				hasError := false

				if startStr != "" {
					parsed, err := time.Parse(dateLayout, startStr)
					if err != nil {
						screen.SetFieldError(fieldFilterStartDate, "используйте формат ГГГГ-ММ-ДД")
						hasError = true
					} else {
						screen.SetFieldError(fieldFilterStartDate, "")
						startDate = &parsed
					}
				} else {
					screen.SetFieldError(fieldFilterStartDate, "")
				}

				if endStr != "" {
					parsed, err := time.Parse(dateLayout, endStr)
					if err != nil {
						screen.SetFieldError(fieldFilterEndDate, "используйте формат ГГГГ-ММ-ДД")
						hasError = true
					} else {
						screen.SetFieldError(fieldFilterEndDate, "")
						endDate = &parsed
					}
				} else {
					screen.SetFieldError(fieldFilterEndDate, "")
				}

				if startDate != nil && endDate != nil && startDate.After(*endDate) {
					screen.SetFieldError(fieldFilterStartDate, "дата начала должна предшествовать окончанию")
					screen.SetFieldError(fieldFilterEndDate, "дата окончания должна следовать после начала")
					hasError = true
				}

				filter := query.NewOperationFilter()

				if accountID != "" {
					filter = filter.ForAccount(domain.ID(accountID))
				}

				if categoryID != "" {
					filter = filter.ForCategory(domain.ID(categoryID))
				}

				if operationType != "" {
					filter = filter.OfType(domain.OperationType(operationType))
				}

				if startDate != nil && endDate != nil {
					filter = filter.Between(*startDate, *endDate)
				} else if startDate != nil {
					filter = filter.From(*startDate)
				} else if endDate != nil {
					filter = filter.To(*endDate)
				}

				if hasError {
					return tui.Result{}
				}

				listCmd := ctx.OperationCommands().List(filter)
				operations, err := listCmd.Execute(ctx.Context())
				if err != nil {
					screen.SetFieldError(fieldFilterStartDate, err.Error())
					return tui.Result{}
				}

				totals, totalsErr := computeTotals(ctx, operations)
				if totalsErr != nil {
					screen.SetFieldError(fieldFilterStartDate, totalsErr.Error())
					return tui.Result{}
				}

				return tui.Result{Push: NewList(filter, operations, accounts, categories, totals)}
			},
		),
		menus.NewPopItem("Назад", "Вернуться в меню операций"),
	}

	screen = menus.NewScreen(
		"Фильтр операций",
		"Настройте параметры и просмотрите список операций.",
		items,
	)

	return screen
}

func computeTotals(ctx tui.ScreenContext, operations []*domain.Operation) (appanalytics.Totals, error) {
	cmdService := ctx.AnalyticsCommands()
	if cmdService == nil {
		return appanalytics.Totals{}, nil
	}

	cmd := cmdService.NetTotals(operations)
	if cmd == nil {
		return appanalytics.Totals{}, nil
	}

	return cmd.Execute(ctx.Context())
}
