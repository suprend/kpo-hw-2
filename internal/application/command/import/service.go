package fileimport

import (
	"context"

	appcommand "kpo-hw-2/internal/application/command"
	appfiles "kpo-hw-2/internal/application/files"
	fileimport "kpo-hw-2/internal/application/files/import"
)

type Service struct {
	importService *fileimport.Service
	decorators    Decorators
}

func NewService(importService *fileimport.Service, decorators Decorators) *Service {
	return &Service{
		importService: importService,
		decorators:    decorators,
	}
}

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

type Decorators struct {
	ListFormats    []appcommand.Decorator[[]appfiles.Format]
	ImportFromPath []appcommand.Decorator[fileimport.Result]
}
