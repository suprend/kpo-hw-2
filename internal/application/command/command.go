package command

import (
	"context"
	"errors"
)

type Command[T any] interface {
	Execute(ctx context.Context) (T, error)
	Name() string
}

type Decorator[T any] interface {
	Wrap(cmd Command[T]) Command[T]
}

func Wrap[T any](cmd Command[T], decorators ...Decorator[T]) Command[T] {
	if cmd == nil {
		return nil
	}

	for _, decorator := range decorators {
		if decorator == nil {
			continue
		}
		cmd = decorator.Wrap(cmd)
	}

	return cmd
}

type Func[T any] struct {
	ExecFn func(context.Context) (T, error)
	NameFn func() string
}

func (f Func[T]) Execute(ctx context.Context) (T, error) {
	if f.ExecFn == nil {
		var zero T
		return zero, errors.New("command.Func: ExecFn is nil")
	}
	return f.ExecFn(ctx)
}

func (f Func[T]) Name() string {
	if f.NameFn == nil {
		return ""
	}
	return f.NameFn()
}
