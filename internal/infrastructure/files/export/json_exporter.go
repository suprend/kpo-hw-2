package fileexport

import (
	"encoding/json"
	"io"

	appfiles "kpo-hw-2/internal/application/files"
	fileexport "kpo-hw-2/internal/application/files/export"
	"kpo-hw-2/internal/domain"
	filesmodel "kpo-hw-2/internal/files/model"
)

type JSONExporter struct{}

func NewJSONExporter() *JSONExporter {
	return &JSONExporter{}
}

func (e *JSONExporter) Format() appfiles.Format {
	return appfiles.Format{
		Key:         "json",
		Title:       "JSON",
		Description: "Экспорт данных в формате JSON.",
		Extension:   "json",
	}
}

func (e *JSONExporter) NewVisitor(writer io.Writer) (fileexport.Visitor, error) {
	return &jsonVisitor{
		writer: writer,
	}, nil
}

var _ fileexport.Exporter = (*JSONExporter)(nil)

type jsonVisitor struct {
	writer  io.Writer
	payload filesmodel.Payload
}

func (v *jsonVisitor) VisitBankAccount(account *domain.BankAccount) error {
	if account == nil {
		return nil
	}
	v.payload.Accounts = append(v.payload.Accounts, filesmodel.Account{
		ID:      account.ID().String(),
		Name:    account.Name(),
		Balance: account.Balance(),
	})
	return nil
}

func (v *jsonVisitor) VisitCategory(category *domain.Category) error {
	if category == nil {
		return nil
	}
	v.payload.Categories = append(v.payload.Categories, filesmodel.Category{
		ID:   category.ID().String(),
		Type: string(category.Type()),
		Name: category.Name(),
	})
	return nil
}

func (v *jsonVisitor) VisitOperation(operation *domain.Operation) error {
	if operation == nil {
		return nil
	}
	v.payload.Operations = append(v.payload.Operations, filesmodel.Operation{
		ID:            operation.ID().String(),
		Type:          string(operation.Type()),
		BankAccountID: operation.BankAccountID().String(),
		CategoryID:    operation.CategoryID().String(),
		Amount:        operation.Amount(),
		Date:          operation.Date(),
		Description:   operation.Description(),
	})
	return nil
}

func (v *jsonVisitor) Finalize() error {
	payload := struct {
		Accounts   []filesmodel.Account   `json:"accounts"`
		Categories []filesmodel.Category  `json:"categories"`
		Operations []filesmodel.Operation `json:"operations"`
	}{
		Accounts:   v.payload.Accounts,
		Categories: v.payload.Categories,
		Operations: v.payload.Operations,
	}

	encoder := json.NewEncoder(v.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(payload)
}

var _ fileexport.Visitor = (*jsonVisitor)(nil)
