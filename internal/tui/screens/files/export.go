package files

import (
	"strings"

	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/menus"
)

const (
	fieldExportFormat = "export_format"
	fieldExportDir    = "export_dir"
	fieldExportName   = "export_name"
)

func newExportScreen(ctx tui.ScreenContext) tui.Screen {
	var screen *menus.Screen

    formats, options, defaultIndex := loadExportFormatsWithOptions(ctx)
	if len(formats) == 0 {
		return emptyFormatsScreen("Экспорт данных")
	}

	defaultFormat := formats[defaultIndex]
	defaultDir, defaultName := defaultExportLocation(defaultFormat)

	items := []menus.MenuItem{
		menus.NewSelectItem(
			fieldExportFormat,
			"Формат",
			"Выберите формат экспорта.",
			options,
			menus.SelectConfig{InitialIndex: defaultIndex},
		),
		menus.NewInputItem(
			fieldExportDir,
			"Папка",
			"Укажите каталог, в котором будет сохранён файл.",
			menus.InputConfig{
				Placeholder: "storage",
				Initial:     defaultDir,
			},
		),
		menus.NewInputItem(
			fieldExportName,
			"Имя файла",
			"Введите имя файла без расширения.",
			menus.InputConfig{
				Placeholder: "data",
				Initial:     defaultName,
			},
		),
		menus.NewActionItem(
			"save",
			"Сохранить",
			"Выполнить экспорт данных.",
			func(context tui.ScreenContext, values menus.Values) tui.Result {
				formatKey := screen.Value(fieldExportFormat)
				format := findFormat(formats, formatKey)

				dir := strings.TrimSpace(values[fieldExportDir])
				if dir == "" {
					dir = defaultDir
					screen.SetValue(fieldExportDir, dir)
				}
				if dir == "" {
					screen.SetFieldError(fieldExportDir, "укажите папку")
					return tui.Result{}
				}

				name := strings.TrimSpace(values[fieldExportName])
				if name == "" {
					name = defaultName
					screen.SetValue(fieldExportName, name)
				}
				if name == "" {
					screen.SetFieldError(fieldExportName, "укажите имя файла")
					return tui.Result{}
				}

				path := exportFilePath(dir, name, format)

				cmd := context.ExportCommands().ExportToPath(formatKey, path)
				if _, err := cmd.Execute(context.Context()); err != nil {
					screen.SetFieldError(fieldExportName, err.Error())
					return tui.Result{}
				}

				successMessage := "Данные сохранены в " + path
				return tui.Result{
					Push: successScreen("Экспорт завершён", successMessage),
				}
			},
		),
		menus.NewPopItem("Назад", "Вернуться к меню файлов"),
	}

	screen = menus.NewScreen(
		"Экспорт данных",
		"Настройте параметры и сохраните данные.",
		items,
	)
	return screen
}
