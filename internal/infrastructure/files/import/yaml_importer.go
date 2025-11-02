package fileimport

import (
	"gopkg.in/yaml.v3"

	appfiles "kpo-hw-2/internal/application/files"
	fileimport "kpo-hw-2/internal/application/files/import"
	filesmodel "kpo-hw-2/internal/files/model"
)

// YAMLImporter парсит данные из YAML, экспортированного приложением.
type YAMLImporter struct{}

// NewYAMLImporter создает экземпляр YAML импортера.
func NewYAMLImporter() *YAMLImporter {
	return &YAMLImporter{}
}

// Format описывает поддерживаемый формат.
func (i *YAMLImporter) Format() appfiles.Format {
	return appfiles.Format{
		Key:         "yaml",
		Title:       "YAML",
		Description: "Импорт данных из YAML-файла.",
		Extension:   "yaml",
	}
}

// Parse преобразует сырые данные в транспортную структуру.
func (i *YAMLImporter) Parse(data []byte) (filesmodel.Payload, error) {
	if len(data) == 0 {
		return filesmodel.Payload{}, nil
	}

	var payload filesmodel.Payload
	if err := yaml.Unmarshal(data, &payload); err != nil {
		return filesmodel.Payload{}, err
	}

	return payload, nil
}

var _ fileimport.Importer = (*YAMLImporter)(nil)
