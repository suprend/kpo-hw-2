package export

import (
	"io"

	"kpo-hw-2/internal/application/files"
)

// Exporter produces visitors that serialize domain aggregates into a specific format.
type Exporter interface {
	// Format returns metadata describing the export format.
	Format() files.Format
	// NewVisitor constructs a visitor bound to the provided writer.
	NewVisitor(writer io.Writer) (Visitor, error)
}
