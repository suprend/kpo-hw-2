package files

// Format описывает формат файла, доступный для импорта/экспорта.
type Format struct {
	Key         string // уникальный идентификатор (например, "json")
	Title       string // отображаемое название (например, "JSON")
	Description string // дополнительное описание
	Extension   string // расширение без точки
}
