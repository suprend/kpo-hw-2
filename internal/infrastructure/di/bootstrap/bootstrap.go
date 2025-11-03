package bootstrap

import (
	"context"
	"fmt"
	"time"

	accountcmd "kpo-hw-2/internal/application/command/account"
	analyticscmd "kpo-hw-2/internal/application/command/analytics"
	categorycmd "kpo-hw-2/internal/application/command/category"
	exportcmd "kpo-hw-2/internal/application/command/export"
	importcmd "kpo-hw-2/internal/application/command/import"
	operationcmd "kpo-hw-2/internal/application/command/operation"
	fileexport "kpo-hw-2/internal/application/files/export"
	fileimport "kpo-hw-2/internal/application/files/import"
	"kpo-hw-2/internal/infrastructure/di"
	"kpo-hw-2/internal/tui"
)

type App struct {
	Model *tui.Model
}

func Build(
	ctx context.Context,
	logFn func(string, time.Duration, error),
	exporters []fileexport.Exporter,
	importers []fileimport.Importer,
) (*App, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	container := di.New()

	if err := di.Provide[func(string, time.Duration, error)](container, logFn); err != nil {
		return nil, fmt.Errorf("bootstrap: provide log function: %w", err)
	}
	if err := di.Provide[[]fileexport.Exporter](container, exporters); err != nil {
		return nil, fmt.Errorf("bootstrap: provide exporters: %w", err)
	}
	if err := di.Provide[[]fileimport.Importer](container, importers); err != nil {
		return nil, fmt.Errorf("bootstrap: provide importers: %w", err)
	}

	if err := registerInfrastructure(container); err != nil {
		return nil, err
	}
	if err := registerDomain(container); err != nil {
		return nil, err
	}
	if err := registerApplication(container); err != nil {
		return nil, err
	}
	if err := registerCommands(container); err != nil {
		return nil, err
	}
	if err := registerUI(container); err != nil {
		return nil, err
	}

	accountCommands, err := di.Resolve[*accountcmd.Service](container)
	if err != nil {
		return nil, fmt.Errorf("bootstrap: resolve account commands: %w", err)
	}
	categoryCommands, err := di.Resolve[*categorycmd.Service](container)
	if err != nil {
		return nil, fmt.Errorf("bootstrap: resolve category commands: %w", err)
	}
	operationCommands, err := di.Resolve[*operationcmd.Service](container)
	if err != nil {
		return nil, fmt.Errorf("bootstrap: resolve operation commands: %w", err)
	}
	fileExportCommands, err := di.Resolve[*exportcmd.Service](container)
	if err != nil {
		return nil, fmt.Errorf("bootstrap: resolve export commands: %w", err)
	}
	fileImportCommands, err := di.Resolve[*importcmd.Service](container)
	if err != nil {
		return nil, fmt.Errorf("bootstrap: resolve import commands: %w", err)
	}
	analyticsCommands, err := di.Resolve[*analyticscmd.Service](container)
	if err != nil {
		return nil, fmt.Errorf("bootstrap: resolve analytics commands: %w", err)
	}

	rootScreen, err := di.Resolve[tui.Screen](container)
	if err != nil {
		return nil, fmt.Errorf("bootstrap: resolve root screen: %w", err)
	}

	model := tui.NewProgram(
		ctx,
		accountCommands,
		categoryCommands,
		operationCommands,
		fileExportCommands,
		fileImportCommands,
		analyticsCommands,
		rootScreen,
	)

	return &App{Model: model}, nil
}
