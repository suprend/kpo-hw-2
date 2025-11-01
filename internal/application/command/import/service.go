package fileimport

import (
	"context"

	appcommand "kpo-hw-2/internal/application/command"
	appfiles "kpo-hw-2/internal/application/files"
	fileimport "kpo-hw-2/internal/application/files/import"
)

// Service строит команды, делегирующие в сервис импорта.
type Service struct {
	importService *fileimport.Service
	decorators    Decorators
}

// NewService регистрирует сервис импорта и декораторы.
func NewService(importService *fileimport.Service, decorators Decorators) *Service {
	return &Service{
		importService: importService,
		decorators:    decorators,
	}
}

// ListFormats возвращает команду для получения доступных форматов импорта.
func (s *Service) ListFormats() appcommand.Command[[]appfiles.Format] {
	base := appcommand.Func[[]appfiles.Format]{
		ExecFn: func(_ context.Context) ([]appfiles.Format, error) {
			if s.importService == nil {
				return nil, nil
			}
			return s.importService.Formats(), nil
		},
		NameFn: func() string { return "import.formats" },
	}
	return appcommand.Wrap(base, s.decorators.ListFormats...)
}

// ImportFromPath возвращает команду, которая загружает данные из файла.
func (s *Service) ImportFromPath(formatKey, path string) appcommand.Command[fileimport.Result] {
	base := appcommand.Func[fileimport.Result]{
		ExecFn: func(_ context.Context) (fileimport.Result, error) {
			if s.importService == nil {
				return fileimport.Result{}, nil
			}
			return s.importService.ImportFromPath(formatKey, path)
		},
		NameFn: func() string { return "import.from_path" },
	}
	return appcommand.Wrap(base, s.decorators.ImportFromPath...)
}

// Decorators группирует опциональные декораторы команд импорта.
type Decorators struct {
	ListFormats    []appcommand.Decorator[[]appfiles.Format]
	ImportFromPath []appcommand.Decorator[fileimport.Result]
}
