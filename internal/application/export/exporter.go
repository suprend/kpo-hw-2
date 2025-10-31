package export

import (
	"io"

	"kpo-hw-2/internal/domain"
)

// Exporter serializes domain entities to desired representation.
type Exporter interface {
	Export(
		writer io.Writer,
		accounts []*domain.BankAccount,
		categories []*domain.Category,
		operations []*domain.Operation,
	) error
}
