package command

import (
	"context"
	"errors"
)

// Command describes a user scenario executed by the application layer.
// The generic result type makes it possible to reuse the same interface
// for commands that return entities or simple acknowledgements.
type Command[T any] interface {
	Execute(ctx context.Context) (T, error)
	Name() string
}

// Decorator wraps commands with cross-cutting behaviour (logging, metrics, etc.).
// Implementations typically return a new Command that delegates to the original one.
type Decorator[T any] interface {
	Wrap(cmd Command[T]) Command[T]
}

// Wrap applies decorators to the command in declaration order.
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

// Func provides a lightweight Command implementation based on callbacks.
type Func[T any] struct {
	ExecFn func(context.Context) (T, error)
	NameFn func() string
}

// Execute delegates to ExecFn or returns an error when callback is missing.
func (f Func[T]) Execute(ctx context.Context) (T, error) {
	if f.ExecFn == nil {
		var zero T
		return zero, errors.New("command.Func: ExecFn is nil")
	}
	return f.ExecFn(ctx)
}

// Name delegates to NameFn and falls back to an empty string.
func (f Func[T]) Name() string {
	if f.NameFn == nil {
		return ""
	}
	return f.NameFn()
}
