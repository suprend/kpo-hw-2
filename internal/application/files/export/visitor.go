package export

import "kpo-hw-2/internal/domain"

// Visitor describes how to process domain aggregates during export.
type Visitor interface {
	VisitBankAccount(*domain.BankAccount) error
	VisitCategory(*domain.Category) error
	VisitOperation(*domain.Operation) error
	Finalize() error
}
