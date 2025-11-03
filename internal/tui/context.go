package tui

import (
	"context"

	accountcmd "kpo-hw-2/internal/application/command/account"
	analyticscmd "kpo-hw-2/internal/application/command/analytics"
	categorycmd "kpo-hw-2/internal/application/command/category"
	exportcmd "kpo-hw-2/internal/application/command/export"
	fileimportcmd "kpo-hw-2/internal/application/command/import"
	operationcmd "kpo-hw-2/internal/application/command/operation"
)

type ScreenContext interface {
	Context() context.Context

	AccountCommands() *accountcmd.Service
	CategoryCommands() *categorycmd.Service
	OperationCommands() *operationcmd.Service
	ExportCommands() *exportcmd.Service
	ImportCommands() *fileimportcmd.Service
	AnalyticsCommands() *analyticscmd.Service
}
