package export

import (
	"io"

	"kpo-hw-2/internal/application/files"
)

type Exporter interface {
	Format() files.Format
	NewVisitor(writer io.Writer) (Visitor, error)
}
