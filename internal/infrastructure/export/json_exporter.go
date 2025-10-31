package export

import (
	"encoding/json"
	"io"
	"time"

	appExport "kpo-hw-2/internal/application/export"
	"kpo-hw-2/internal/domain"
)

// JSONExporter writes domain entities into JSON representation.
type JSONExporter struct{}

// NewJSONExporter constructs exporter instance.
func NewJSONExporter() *JSONExporter {
	return &JSONExporter{}
}

// Export serializes provided entities into writer as JSON.
func (e *JSONExporter) Export(
	writer io.Writer,
	accounts []*domain.BankAccount,
	categories []*domain.Category,
	operations []*domain.Operation,
) error {
	payload := struct {
		Accounts   []accountDTO   `json:"accounts"`
		Categories []categoryDTO  `json:"categories"`
		Operations []operationDTO `json:"operations"`
	}{
		Accounts:   make([]accountDTO, 0, len(accounts)),
		Categories: make([]categoryDTO, 0, len(categories)),
		Operations: make([]operationDTO, 0, len(operations)),
	}

	for _, acc := range accounts {
		payload.Accounts = append(payload.Accounts, accountDTO{
			ID:      acc.ID().String(),
			Name:    acc.Name(),
			Balance: acc.Balance(),
		})
	}

	for _, cat := range categories {
		payload.Categories = append(payload.Categories, categoryDTO{
			ID:   cat.ID().String(),
			Type: string(cat.Type()),
			Name: cat.Name(),
		})
	}

	for _, op := range operations {
		payload.Operations = append(payload.Operations, operationDTO{
			ID:            op.ID().String(),
			Type:          string(op.Type()),
			BankAccountID: op.BankAccountID().String(),
			CategoryID:    op.CategoryID().String(),
			Amount:        op.Amount(),
			Date:          op.Date(),
			Description:   op.Description(),
		})
	}

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(payload)
}

// compile-time check
var _ appExport.Exporter = (*JSONExporter)(nil)

type accountDTO struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Balance int64  `json:"balance"`
}

type categoryDTO struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type operationDTO struct {
	ID            string    `json:"id"`
	Type          string    `json:"type"`
	BankAccountID string    `json:"bank_account_id"`
	CategoryID    string    `json:"category_id"`
	Amount        int64     `json:"amount"`
	Date          time.Time `json:"date"`
	Description   string    `json:"description,omitempty"`
}
