package fileexport

import (
	"encoding/csv"
	"io"
	"strconv"
	"time"

	appfiles "kpo-hw-2/internal/application/files"
	fileexport "kpo-hw-2/internal/application/files/export"
	"kpo-hw-2/internal/domain"
	filesmodel "kpo-hw-2/internal/files/model"
)

type CSVExporter struct{}

func NewCSVExporter() *CSVExporter {
	return &CSVExporter{}
}

func (e *CSVExporter) Format() appfiles.Format {
	return appfiles.Format{
		Key:         "csv",
		Title:       "CSV",
		Description: "Экспорт данных в формате CSV.",
		Extension:   "csv",
	}
}

func (e *CSVExporter) NewVisitor(writer io.Writer) (fileexport.Visitor, error) {
	return &csvVisitor{
		writer: csv.NewWriter(writer),
	}, nil
}

var _ fileexport.Exporter = (*CSVExporter)(nil)

type csvVisitor struct {
	writer  *csv.Writer
	payload filesmodel.Payload
}

func (v *csvVisitor) VisitBankAccount(account *domain.BankAccount) error {
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

func (v *csvVisitor) VisitCategory(category *domain.Category) error {
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

func (v *csvVisitor) VisitOperation(operation *domain.Operation) error {
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

func (v *csvVisitor) Finalize() error {
	if v.writer == nil {
		return nil
	}

	if err := v.writer.Write([]string{
		"entity",
		"id",
		"name",
		"type",
		"balance",
		"bank_account_id",
		"category_id",
		"amount",
		"date",
		"description",
	}); err != nil {
		return err
	}

	for _, account := range v.payload.Accounts {
		if err := v.writer.Write([]string{
			"account",
			account.ID,
			account.Name,
			"",
			strconv.FormatInt(account.Balance, 10),
			"",
			"",
			"",
			"",
			"",
		}); err != nil {
			return err
		}
	}

	for _, category := range v.payload.Categories {
		if err := v.writer.Write([]string{
			"category",
			category.ID,
			category.Name,
			category.Type,
			"",
			"",
			"",
			"",
			"",
			"",
		}); err != nil {
			return err
		}
	}

	for _, operation := range v.payload.Operations {
		dateValue := ""
		if !operation.Date.IsZero() {
			dateValue = operation.Date.Format(time.RFC3339)
		}

		if err := v.writer.Write([]string{
			"operation",
			operation.ID,
			"",
			operation.Type,
			"",
			operation.BankAccountID,
			operation.CategoryID,
			strconv.FormatInt(operation.Amount, 10),
			dateValue,
			operation.Description,
		}); err != nil {
			return err
		}
	}

	v.writer.Flush()
	return v.writer.Error()
}

var _ fileexport.Visitor = (*csvVisitor)(nil)
