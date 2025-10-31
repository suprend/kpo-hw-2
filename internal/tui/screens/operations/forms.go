package operations

import (
	"fmt"

	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/tui/menus"
)

type selectData struct {
	options   []menus.SelectOption
	indexByID map[string]int
	typeByID  map[string]domain.OperationType
}

func buildAccountSelectData(accounts []*domain.BankAccount) selectData {
	data := selectData{
		options:   make([]menus.SelectOption, 0, len(accounts)),
		indexByID: make(map[string]int, len(accounts)),
	}

	for idx, account := range accounts {
		id := account.ID().String()
		data.options = append(data.options, menus.SelectOption{
			Label: fmt.Sprintf("%s (%d)", account.Name(), account.Balance()),
			Value: id,
		})
		data.indexByID[id] = idx
	}

	return data
}

func buildCategorySelectData(categories []*domain.Category) selectData {
	data := selectData{
		options:   make([]menus.SelectOption, 0, len(categories)),
		indexByID: make(map[string]int, len(categories)),
		typeByID:  make(map[string]domain.OperationType, len(categories)),
	}

	for idx, category := range categories {
		id := category.ID().String()
		data.options = append(data.options, menus.SelectOption{
			Label: fmt.Sprintf("%s (%s)", category.Name(), readableType(category.Type())),
			Value: id,
		})
		data.indexByID[id] = idx
		data.typeByID[id] = category.Type()
	}

	return data
}
