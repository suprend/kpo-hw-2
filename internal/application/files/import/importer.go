package fileimport

import (
	"kpo-hw-2/internal/application/files"
	filesmodel "kpo-hw-2/internal/files/model"
)

type Importer interface {
	Format() files.Format
	Parse(data []byte) (filesmodel.Payload, error)
}
