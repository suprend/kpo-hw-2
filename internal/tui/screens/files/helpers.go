package files

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	appfiles "kpo-hw-2/internal/application/files"
	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/menus"
)

func loadExportFormatsWithOptions(ctx tui.ScreenContext) ([]appfiles.Format, []menus.SelectOption, int) {
	return loadFormats(ctx.ExportCommands().ListFormats(), ctx)
}

func loadImportFormatsWithOptions(ctx tui.ScreenContext) ([]appfiles.Format, []menus.SelectOption, int) {
	if ctx.ImportCommands() == nil {
		return nil, nil, 0
	}
	return loadFormats(ctx.ImportCommands().ListFormats(), ctx)
}

func loadFormats(cmd interface {
	Execute(context.Context) ([]appfiles.Format, error)
}, ctx tui.ScreenContext) ([]appfiles.Format, []menus.SelectOption, int) {
	formats, err := cmd.Execute(ctx.Context())
	if err != nil || len(formats) == 0 {
		return nil, nil, 0
	}

	options := make([]menus.SelectOption, len(formats))
	for i, format := range formats {
		label := format.Title
		if label == "" {
			label = format.Key
		}
		options[i] = menus.SelectOption{Label: label, Value: format.Key}
	}

	return formats, options, 0
}

func defaultExportLocation(format appfiles.Format) (string, string) {
	dir := "storage"
	if abs, err := filepath.Abs(dir); err == nil {
		dir = abs
	}

	return dir, "data"
}

func exportFilePath(dir, name string, format appfiles.Format) string {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		dir = "."
	}
	name = strings.TrimSpace(name)
	if name == "" {
		name = "data"
	}

	ext := strings.TrimLeft(format.Extension, ".")
	if ext == "" {
		ext = "json"
	}

	filename := fmt.Sprintf("%s.%s", name, ext)
	return filepath.Join(dir, filename)
}

func findFormat(formats []appfiles.Format, key string) appfiles.Format {
	for _, format := range formats {
		if format.Key == key {
			return format
		}
	}
	return appfiles.Format{}
}

func successScreen(title, message string) tui.Screen {
	return menus.NewScreen(
		title,
		message,
		[]menus.MenuItem{menus.NewPopItem("Ок", "Вернуться к меню файлов")},
	)
}

func emptyFormatsScreen(title string) tui.Screen {
	return menus.NewScreen(
		title,
		"Нет доступных форматов.",
		[]menus.MenuItem{menus.NewPopItem("Назад", "Вернуться к меню файлов")},
	)
}
