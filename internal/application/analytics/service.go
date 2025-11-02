package analytics

import (
	"fmt"

	"kpo-hw-2/internal/domain"
)

// Totals объединяет итоговые показатели по операциям.
type Totals struct {
	Income  int64
	Expense int64
	Delta   int64
}

// Service описывает минимальный набор аналитики.
type Service interface {
	NetTotals(operations []*domain.Operation) (Totals, error)
}

type service struct{}

// NewService возвращает готовую реализацию Service.
func NewService() Service {
	return service{}
}

func (service) NetTotals(operations []*domain.Operation) (Totals, error) {
	var totals Totals

	for _, op := range operations {
		if op == nil {
			continue
		}

		switch op.Type() {
		case domain.OperationTypeIncome:
			totals.Income += op.Amount()
		case domain.OperationTypeExpense:
			totals.Expense += op.Amount()
		default:
			return Totals{}, fmt.Errorf("analytics: unsupported operation type %q", op.Type())
		}
	}

	totals.Delta = totals.Income - totals.Expense
	return totals, nil
}
