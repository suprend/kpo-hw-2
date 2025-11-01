package export

import (
	"context"
	"io"

	appcommand "kpo-hw-2/internal/application/command"
	appfiles "kpo-hw-2/internal/application/files"
	fileexport "kpo-hw-2/internal/application/files/export"
)

// Service constructs commands that delegate to the export application service.
type Service struct {
	exportService *fileexport.Service
	decorators    Decorators
}

// NewService wires export service with configured decorators.
func NewService(exportService *fileexport.Service, decorators Decorators) *Service {
	return &Service{
		exportService: exportService,
		decorators:    decorators,
	}
}

// ListFormats returns a command that fetches available export formats.
func (s *Service) ListFormats() appcommand.Command[[]appfiles.Format] {
	base := appcommand.Func[[]appfiles.Format]{
		ExecFn: func(_ context.Context) ([]appfiles.Format, error) {
			if s.exportService == nil {
				return nil, nil
			}
			return s.exportService.Formats(), nil
		},
		NameFn: func() string { return "export.formats" },
	}
	return appcommand.Wrap(base, s.decorators.ListFormats...)
}

// ExportToPath returns a command that exports data into file under provided path.
func (s *Service) ExportToPath(formatKey, destination string) appcommand.Command[appcommand.NoResult] {
	base := appcommand.Func[appcommand.NoResult]{
		ExecFn: func(_ context.Context) (appcommand.NoResult, error) {
			if s.exportService == nil {
				return appcommand.NoResult{}, nil
			}
			err := s.exportService.ExportToPath(formatKey, destination)
			return appcommand.NoResult{}, err
		},
		NameFn: func() string { return "export.to_path" },
	}
	return appcommand.Wrap(base, s.decorators.ExportToPath...)
}

// Export returns a command that writes export data into provided writer.
func (s *Service) Export(formatKey string, writer io.Writer) appcommand.Command[appcommand.NoResult] {
	base := appcommand.Func[appcommand.NoResult]{
		ExecFn: func(_ context.Context) (appcommand.NoResult, error) {
			if s.exportService == nil {
				return appcommand.NoResult{}, nil
			}
			err := s.exportService.Export(formatKey, writer)
			return appcommand.NoResult{}, err
		},
		NameFn: func() string { return "export.write" },
	}
	return appcommand.Wrap(base, s.decorators.Export...)
}

// Decorators groups optional decorators for export commands.
type Decorators struct {
	ListFormats  []appcommand.Decorator[[]appfiles.Format]
	ExportToPath []appcommand.Decorator[appcommand.NoResult]
	Export       []appcommand.Decorator[appcommand.NoResult]
}
