package fileimport

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"kpo-hw-2/internal/application/facade"
	"kpo-hw-2/internal/application/files"
	"kpo-hw-2/internal/domain"
	filesmodel "kpo-hw-2/internal/files/model"
)

var (
	ErrUnknownFormat = errors.New("import: unknown format")
	ErrInvalidSource = errors.New("import: invalid source")
	ErrInvalidPath   = errors.New("import: invalid source path")
)

type Result struct {
	CreatedAccounts   int
	CreatedCategories int
	CreatedOperations int

	SkippedAccounts   int
	SkippedCategories int
	SkippedOperations int
}

type Service struct {
	accounts   facade.AccountFacade
	categories facade.CategoryFacade
	operations facade.OperationFacade

	importers map[string]Importer
	order     []files.Format
}

func NewService(
	accountFacade facade.AccountFacade,
	categoryFacade facade.CategoryFacade,
	operationFacade facade.OperationFacade,
	importers []Importer,
) *Service {
	registry := make(map[string]Importer)
	order := make([]files.Format, 0, len(importers))

	for _, imp := range importers {
		if imp == nil {
			continue
		}
		format := imp.Format()
		if format.Key == "" {
			continue
		}
		if _, exists := registry[format.Key]; exists {
			continue
		}
		registry[format.Key] = imp
		order = append(order, format)
	}

	return &Service{
		accounts:   accountFacade,
		categories: categoryFacade,
		operations: operationFacade,
		importers:  registry,
		order:      order,
	}
}

func (s *Service) Formats() []files.Format {
	out := make([]files.Format, len(s.order))
	copy(out, s.order)
	return out
}

func (s *Service) FormatFor(key string) (files.Format, bool) {
	imp, ok := s.importers[key]
	if !ok {
		return files.Format{}, false
	}
	return imp.Format(), true
}

func (s *Service) ImportFromPath(formatKey, path string) (Result, error) {
	if strings.TrimSpace(path) == "" {
		return Result{}, ErrInvalidPath
	}

	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return Result{}, err
	}
	defer file.Close()

	return s.Import(formatKey, file)
}

func (s *Service) Import(formatKey string, reader io.Reader) (Result, error) {
	if reader == nil {
		return Result{}, ErrInvalidSource
	}

	importer, ok := s.importers[formatKey]
	if !ok {
		return Result{}, ErrUnknownFormat
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return Result{}, err
	}

	payload, err := importer.Parse(data)
	if err != nil {
		return Result{}, err
	}

	return s.applyPayload(payload)
}

func (s *Service) applyPayload(payload filesmodel.Payload) (Result, error) {
	var result Result

	accountIDs := make(map[string]domain.ID)
	categoryIDs := make(map[string]domain.ID)

	if s.accounts != nil {
		for _, dto := range payload.Accounts {
			name := strings.TrimSpace(dto.Name)
			id := domain.ID(strings.TrimSpace(dto.ID))
			if id == "" {
				result.SkippedAccounts++
				continue
			}

			account, err := s.accounts.CreateAccountWithID(id, name, dto.Balance)
			if err != nil {
				if errors.Is(err, domain.ErrAlreadyExists) {
					accountIDs[dto.ID] = id
					result.SkippedAccounts++
					continue
				}

				result.SkippedAccounts++
				continue
			}

			accountIDs[dto.ID] = account.ID()
			result.CreatedAccounts++
		}
	}

	if s.categories != nil {
		for _, dto := range payload.Categories {
			name := strings.TrimSpace(dto.Name)
			typ := domain.OperationType(strings.ToLower(strings.TrimSpace(dto.Type)))
			id := domain.ID(strings.TrimSpace(dto.ID))
			if id == "" {
				result.SkippedCategories++
				continue
			}

			category, err := s.categories.CreateCategoryWithID(id, name, typ)
			if err != nil {
				if errors.Is(err, domain.ErrAlreadyExists) {
					categoryIDs[dto.ID] = id
					result.SkippedCategories++
					continue
				}

				result.SkippedCategories++
				continue
			}

			categoryIDs[dto.ID] = category.ID()
			result.CreatedCategories++
		}
	}

	if s.operations != nil {
		for _, dto := range payload.Operations {
			typ := domain.OperationType(strings.ToLower(strings.TrimSpace(dto.Type)))
			id := domain.ID(strings.TrimSpace(dto.ID))
			if id == "" {
				result.SkippedOperations++
				continue
			}

			accountID, ok := accountIDs[dto.BankAccountID]
			if !ok {
				result.SkippedOperations++
				continue
			}

			categoryID, ok := categoryIDs[dto.CategoryID]
			if !ok {
				result.SkippedOperations++
				continue
			}

			if _, err := s.operations.CreateOperationWithoutBalance(
				id,
				typ,
				accountID,
				categoryID,
				dto.Amount,
				dto.Date,
				strings.TrimSpace(dto.Description),
			); err != nil {
				if errors.Is(err, domain.ErrAlreadyExists) {
					result.SkippedOperations++
					continue
				}

				result.SkippedOperations++
				continue
			}

			result.CreatedOperations++
		}
	}

	return result, nil
}
