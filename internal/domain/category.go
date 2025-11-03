package domain

import "strings"

type Category struct {
	id   ID
	typ  OperationType
	name string
}

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

func (c *Category) ID() ID { return c.id }

func (c *Category) Type() OperationType { return c.typ }

func (c *Category) Name() string { return c.name }
