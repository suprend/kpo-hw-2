package command

import "context"

// NoResult is a convenience type for commands that do not return a value.
type NoResult struct{}

// Factory stores cross-cutting decorators to produce wrapped commands.
type Factory[T any] struct {
	decorators []Decorator[T]
}

// NewFactory constructs a factory with predefined decorators.
func NewFactory[T any](decorators ...Decorator[T]) Factory[T] {
	return Factory[T]{decorators: decorators}
}

// WithDecorators returns a new factory with extra decorators appended.
func (f Factory[T]) WithDecorators(decorators ...Decorator[T]) Factory[T] {
	next := make([]Decorator[T], 0, len(f.decorators)+len(decorators))
	next = append(next, f.decorators...)
	next = append(next, decorators...)
	return Factory[T]{decorators: next}
}

// Wrap applies factory decorators to the provided command.
func (f Factory[T]) Wrap(cmd Command[T]) Command[T] {
	return Wrap(cmd, f.decorators...)
}

// Func creates a command from callbacks and applies factory decorators.
func (f Factory[T]) Func(
	name string,
	exec func(context.Context) (T, error),
) Command[T] {
	cmd := Func[T]{
		ExecFn: exec,
		NameFn: func() string { return name },
	}
	return f.Wrap(cmd)
}
