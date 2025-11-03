package fileexport

import (
	"io"

	"gopkg.in/yaml.v3"

	appfiles "kpo-hw-2/internal/application/files"
	fileexport "kpo-hw-2/internal/application/files/export"
	"kpo-hw-2/internal/domain"
	filesmodel "kpo-hw-2/internal/files/model"
)

type YAMLExporter struct{}

func NewYAMLExporter() *YAMLExporter {
	return &YAMLExporter{}
}

func (e *YAMLExporter) Format() appfiles.Format {
	return appfiles.Format{
		Key:         "yaml",
		Title:       "YAML",
		Description: "Экспорт данных в формате YAML.",
		Extension:   "yaml",
	}
}

func (e *YAMLExporter) NewVisitor(writer io.Writer) (fileexport.Visitor, error) {
	return &yamlVisitor{
		writer: writer,
	}, nil
}

var _ fileexport.Exporter = (*YAMLExporter)(nil)

type yamlVisitor struct {
	writer  io.Writer
	payload filesmodel.Payload
}

func (v *yamlVisitor) VisitBankAccount(account *domain.BankAccount) error {
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

func (v *yamlVisitor) VisitCategory(category *domain.Category) error {
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

func (v *yamlVisitor) VisitOperation(operation *domain.Operation) error {
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

func (v *yamlVisitor) Finalize() error {
	payload := struct {
		Accounts   []filesmodel.Account   `yaml:"accounts"`
		Categories []filesmodel.Category  `yaml:"categories"`
		Operations []filesmodel.Operation `yaml:"operations"`
	}{
		Accounts:   v.payload.Accounts,
		Categories: v.payload.Categories,
		Operations: v.payload.Operations,
	}

	encoder := yaml.NewEncoder(v.writer)
	encoder.SetIndent(2)
	if err := encoder.Encode(payload); err != nil {
		_ = encoder.Close()
		return err
	}
	return encoder.Close()
}

var _ fileexport.Visitor = (*yamlVisitor)(nil)
