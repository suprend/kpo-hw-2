package fileimport

import (
	"encoding/json"

	appfiles "kpo-hw-2/internal/application/files"
	fileimport "kpo-hw-2/internal/application/files/import"
	filesmodel "kpo-hw-2/internal/files/model"
)

// JSONImporter парсит данные из JSON, экспортированного приложением.
type JSONImporter struct{}

// NewJSONImporter создает экземпляр JSON импортера.
func NewJSONImporter() *JSONImporter {
	return &JSONImporter{}
}

// Format описывает поддерживаемый формат.
func (i *JSONImporter) Format() appfiles.Format {
	return appfiles.Format{
		Key:         "json",
		Title:       "JSON",
		Description: "Импорт данных из JSON-файла.",
		Extension:   "json",
	}
}

// Parse преобразует сырые данные в транспортную структуру.
func (i *JSONImporter) Parse(data []byte) (filesmodel.Payload, error) {
	if len(data) == 0 {
		return filesmodel.Payload{}, nil
	}

	var dto payloadDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return filesmodel.Payload{}, err
	}

	return filesmodel.Payload{
		Accounts:   dto.Accounts,
		Categories: dto.Categories,
		Operations: dto.Operations,
	}, nil
}

type payloadDTO struct {
	Accounts   []filesmodel.Account   `json:"accounts"`
	Categories []filesmodel.Category  `json:"categories"`
	Operations []filesmodel.Operation `json:"operations"`
}

var _ fileimport.Importer = (*JSONImporter)(nil)
