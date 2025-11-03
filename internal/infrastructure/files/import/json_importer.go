package fileimport

import (
	"encoding/json"

	appfiles "kpo-hw-2/internal/application/files"
	fileimport "kpo-hw-2/internal/application/files/import"
	filesmodel "kpo-hw-2/internal/files/model"
)

type JSONImporter struct{}

func NewJSONImporter() *JSONImporter {
	return &JSONImporter{}
}

func (i *JSONImporter) Format() appfiles.Format {
	return appfiles.Format{
		Key:         "json",
		Title:       "JSON",
		Description: "Импорт данных из JSON-файла.",
		Extension:   "json",
	}
}

func (i *JSONImporter) Parse(data []byte) (filesmodel.Payload, error) {
	if len(data) == 0 {
		return filesmodel.Payload{}, nil
	}

	var payload filesmodel.Payload
	if err := json.Unmarshal(data, &payload); err != nil {
		return filesmodel.Payload{}, err
	}

	return payload, nil
}

var _ fileimport.Importer = (*JSONImporter)(nil)
