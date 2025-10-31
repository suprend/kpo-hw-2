package domain

import "strings"

// Category represents a classification for operations.
type Category struct {
	id   ID
	typ  OperationType
	name string
}

// NewCategory validates and constructs a category aggregate.
func NewCategory(id ID, typ OperationType, name string) (*Category, error) {
	if id == "" {
		return nil, ErrInvalidCategory
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidCategory
	}

	switch typ {
	case OperationTypeIncome, OperationTypeExpense:
	default:
		return nil, ErrInvalidCategory
	}

	return &Category{
		id:   id,
		typ:  typ,
		name: name,
	}, nil
}

// ID returns the category identifier.
func (c *Category) ID() ID { return c.id }

// Type returns the category type (income/expense).
func (c *Category) Type() OperationType { return c.typ }

// Name returns the category descriptive name.
func (c *Category) Name() string { return c.name }
