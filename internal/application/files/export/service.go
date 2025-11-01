package export

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"kpo-hw-2/internal/application/files"
	"kpo-hw-2/internal/domain/query"
	"kpo-hw-2/internal/domain/repository"
)

var (
	ErrUnknownFormat = errors.New("export: unknown format")
	ErrInvalidWriter = errors.New("export: invalid writer")
	ErrInvalidPath   = errors.New("export: invalid destination path")
)

// Service orchestrates exporting aggregates using registered exporters.
type Service struct {
	accounts   repository.AccountRepository
	categories repository.CategoryRepository
	operations repository.OperationRepository

	exporters map[string]Exporter
	order     []files.Format
}

// NewService wires repository dependencies and registers exporters.
func NewService(
	accountRepo repository.AccountRepository,
	categoryRepo repository.CategoryRepository,
	operationRepo repository.OperationRepository,
	exporters []Exporter,
) *Service {
	registry := make(map[string]Exporter)
	order := make([]files.Format, 0, len(exporters))

	for _, exp := range exporters {
		if exp == nil {
			continue
		}
		format := exp.Format()
		if format.Key == "" {
			continue
		}
		if _, exists := registry[format.Key]; exists {
			continue
		}
		registry[format.Key] = exp
		order = append(order, format)
	}

	return &Service{
		accounts:   accountRepo,
		categories: categoryRepo,
		operations: operationRepo,
		exporters:  registry,
		order:      order,
	}
}

// Formats returns registered formats in deterministic order.
func (s *Service) Formats() []files.Format {
	out := make([]files.Format, len(s.order))
	copy(out, s.order)
	return out
}

// FormatFor returns format metadata for requested key.
func (s *Service) FormatFor(key string) (files.Format, bool) {
	exp, ok := s.exporters[key]
	if !ok {
		return files.Format{}, false
	}
	return exp.Format(), true
}

// Export serializes aggregates into provided writer using selected format.
func (s *Service) Export(formatKey string, writer io.Writer) error {
	if writer == nil {
		return ErrInvalidWriter
	}

	exp, ok := s.exporters[formatKey]
	if !ok {
		return ErrUnknownFormat
	}

	visitor, err := exp.NewVisitor(writer)
	if err != nil {
		return err
	}

	if err := s.exportAccounts(visitor); err != nil {
		return err
	}
	if err := s.exportCategories(visitor); err != nil {
		return err
	}
	if err := s.exportOperations(visitor); err != nil {
		return err
	}

	return visitor.Finalize()
}

// ExportToPath ensures destination exists, creates file and exports into it.
func (s *Service) ExportToPath(formatKey, path string) (err error) {
	if strings.TrimSpace(path) == "" {
		return ErrInvalidPath
	}

	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if mkErr := os.MkdirAll(dir, 0o755); mkErr != nil {
			return mkErr
		}
	}

	file, createErr := os.Create(path)
	if createErr != nil {
		return createErr
	}
	defer func() {
		if cerr := file.Close(); err == nil {
			err = cerr
		}
	}()

	err = s.Export(formatKey, file)
	return err
}

func (s *Service) exportAccounts(visitor Visitor) error {
	if s.accounts == nil {
		return nil
	}

	accounts, err := s.accounts.List()
	if err != nil {
		return err
	}
	for _, account := range accounts {
		if account == nil {
			continue
		}
		if err := visitor.VisitBankAccount(account); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) exportCategories(visitor Visitor) error {
	if s.categories == nil {
		return nil
	}

	categories, err := s.categories.ListAll()
	if err != nil {
		return err
	}
	for _, category := range categories {
		if category == nil {
			continue
		}
		if err := visitor.VisitCategory(category); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) exportOperations(visitor Visitor) error {
	if s.operations == nil {
		return nil
	}

	operations, err := s.operations.ListByFilter(query.NewOperationFilter())
	if err != nil {
		return err
	}
	for _, operation := range operations {
		if operation == nil {
			continue
		}
		if err := visitor.VisitOperation(operation); err != nil {
			return err
		}
	}
	return nil
}
