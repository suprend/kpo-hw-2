package operations

import (
	"context"
	"fmt"
	"strings"

	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/domain/query"
	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/menus"
)

// NewList builds a read-only list of operations based on applied filter.
func NewList(filter query.OperationFilter, operations []*domain.Operation, accounts []*domain.BankAccount, categories []*domain.Category) tui.Screen {
	accountNames := make(map[domain.ID]string, len(accounts))
	for _, acc := range accounts {
		accountNames[acc.ID()] = acc.Name()
	}

	categoryNames := make(map[domain.ID]string, len(categories))
	for _, cat := range categories {
		categoryNames[cat.ID()] = cat.Name()
	}

	items := make([]menus.MenuItem, 0, len(operations)+1)
	for _, op := range operations {
		op := op

		title := fmt.Sprintf("%s • %s", op.Date().Format(dateLayout), operationTitle(op))
		description := buildOperationDescription(op, accountNames, categoryNames)

		items = append(items, menus.NewActionItem(
			op.ID().String(),
			title,
			description,
			func(ctx tui.ScreenContext, _ menus.Values) tui.Result {
				getCmd := ctx.OperationCommands().Get(op.ID())
				operation, err := getCmd.Execute(context.Background())
				if err != nil {
					return tui.Result{}
				}

				accountsCmd := ctx.AccountCommands().List()
				accounts, err := accountsCmd.Execute(context.Background())
				if err != nil {
					return tui.Result{}
				}

				categoriesCmd := ctx.CategoryCommands().List("")
				categories, err := categoriesCmd.Execute(context.Background())
				if err != nil {
					return tui.Result{}
				}

				return tui.Result{Push: NewEdit(operation, accounts, categories)}
			},
		))
	}

	items = append(items, menus.NewPopItem("Назад", "Вернуться к фильтрам"))

	intro := buildFilterIntro(filter, accountNames, categoryNames)

	return menus.NewScreen(
		"Список операций",
		intro,
		items,
	).WithEmptyMessage("Операции не найдены для выбранных параметров.")
}

func operationTitle(op *domain.Operation) string {
	desc := strings.TrimSpace(op.Description())
	if desc != "" {
		return desc
	}
	switch op.Type() {
	case domain.OperationTypeIncome:
		return "Поступление"
	case domain.OperationTypeExpense:
		return "Списание"
	default:
		return "Операция"
	}
}

func buildOperationDescription(op *domain.Operation, accountNames map[domain.ID]string, categoryNames map[domain.ID]string) string {
	amount := op.Amount()
	sign := "+"
	if op.Type() == domain.OperationTypeExpense {
		sign = "-"
	}

	accountName := accountNames[op.BankAccountID()]
	if accountName == "" {
		accountName = op.BankAccountID().String()
	}

	categoryName := categoryNames[op.CategoryID()]
	if categoryName == "" {
		categoryName = op.CategoryID().String()
	}

	return fmt.Sprintf("Сумма: %s%d • Счёт: %s • Категория: %s", sign, amount, accountName, categoryName)
}

func buildFilterIntro(filter query.OperationFilter, accountNames map[domain.ID]string, categoryNames map[domain.ID]string) string {
	var parts []string

	if accID := filter.AccountID(); accID != "" {
		if name, ok := accountNames[accID]; ok {
			parts = append(parts, fmt.Sprintf("Счёт: %s", name))
		} else {
			parts = append(parts, fmt.Sprintf("Счёт: %s", accID.String()))
		}
	} else {
		parts = append(parts, "Все счета")
	}

	if catID := filter.CategoryID(); catID != "" {
		if name, ok := categoryNames[catID]; ok {
			parts = append(parts, fmt.Sprintf("Категория: %s", name))
		} else {
			parts = append(parts, fmt.Sprintf("Категория: %s", catID.String()))
		}
	} else {
		parts = append(parts, "Все категории")
	}

	if from, to := filter.Period(); from != nil || to != nil {
		var periodParts []string
		if from != nil {
			periodParts = append(periodParts, fmt.Sprintf("с %s", from.Format(dateLayout)))
		}
		if to != nil {
			periodParts = append(periodParts, fmt.Sprintf("по %s", to.Format(dateLayout)))
		}
		parts = append(parts, strings.Join(periodParts, " "))
	} else {
		parts = append(parts, "Без ограничения по датам")
	}

	switch typ := filter.Type(); typ {
	case domain.OperationTypeIncome:
		parts = append(parts, "Тип: доход")
	case domain.OperationTypeExpense:
		parts = append(parts, "Тип: расход")
	default:
		parts = append(parts, "Тип: все операции")
	}

	return strings.Join(parts, " • ")
}
