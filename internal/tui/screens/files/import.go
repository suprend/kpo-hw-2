package files

import (
	"fmt"
	"strings"

	fileimport "kpo-hw-2/internal/application/files/import"
	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/menus"
)

const (
	fieldImportFormat = "import_format"
	fieldImportDir    = "import_dir"
	fieldImportName   = "import_name"
)

func newImportScreen(ctx tui.ScreenContext) tui.Screen {
	if ctx.ImportCommands() == nil {
		return emptyFormatsScreen("Импорт данных")
	}

	var screen *menus.Screen

	formats, options, defaultIndex := loadImportFormatsWithOptions(ctx)
	if len(formats) == 0 {
		return emptyFormatsScreen("Импорт данных")
	}

	defaultFormat := formats[defaultIndex]
	defaultDir, defaultName := defaultExportLocation(defaultFormat)

	items := []menus.MenuItem{
		menus.NewSelectItem(
			fieldImportFormat,
			"Формат",
			"Выберите формат импорта.",
			options,
			menus.SelectConfig{InitialIndex: defaultIndex},
		),
		menus.NewInputItem(
			fieldImportDir,
			"Папка",
			"Укажите каталог, где находится файл.",
			menus.InputConfig{
				Placeholder: "storage",
				Initial:     defaultDir,
			},
		),
		menus.NewInputItem(
			fieldImportName,
			"Имя файла",
			"Введите имя файла без расширения.",
			menus.InputConfig{
				Placeholder: "data",
				Initial:     defaultName,
			},
		),
		menus.NewActionItem(
			"load",
			"Загрузить",
			"Импортировать данные из файла.",
			func(context tui.ScreenContext, values menus.Values) tui.Result {
				formatKey := screen.Value(fieldImportFormat)
				format := findFormat(formats, formatKey)

				dir := strings.TrimSpace(values[fieldImportDir])
				if dir == "" {
					dir = defaultDir
					screen.SetValue(fieldImportDir, dir)
				}
				if dir == "" {
					screen.SetFieldError(fieldImportDir, "укажите папку")
					return tui.Result{}
				}

				name := strings.TrimSpace(values[fieldImportName])
				if name == "" {
					name = defaultName
					screen.SetValue(fieldImportName, name)
				}
				if name == "" {
					screen.SetFieldError(fieldImportName, "укажите имя файла")
					return tui.Result{}
				}

				path := exportFilePath(dir, name, format)

				cmd := context.ImportCommands().ImportFromPath(formatKey, path)
				result, err := cmd.Execute(context.Context())
				if err != nil {
					screen.SetFieldError(fieldImportName, err.Error())
					return tui.Result{}
				}

				message := formatImportResult(path, result)
				return tui.Result{
					Push: successScreen("Импорт завершён", message),
				}
			},
		),
		menus.NewPopItem("Назад", "Вернуться к меню файлов"),
	}

	screen = menus.NewScreen(
		"Импорт данных",
		"Укажите параметры и загрузите данные.",
		items,
	)
	return screen
}

func formatImportResult(path string, result fileimport.Result) string {
	return fmt.Sprintf(
		"Данные загружены из %s.\nСоздано: %d счетов, %d категорий, %d операций.\nПропущено: %d/%d/%d.",
		path,
		result.CreatedAccounts,
		result.CreatedCategories,
		result.CreatedOperations,
		result.SkippedAccounts,
		result.SkippedCategories,
		result.SkippedOperations,
	)
}
