package command

import "context"

type NoResult struct{}

type Factory[T any] struct {
	decorators []Decorator[T]
}

func NewFactory[T any](decorators ...Decorator[T]) Factory[T] {
	return Factory[T]{decorators: decorators}
}

func (f Factory[T]) WithDecorators(decorators ...Decorator[T]) Factory[T] {
	next := make([]Decorator[T], 0, len(f.decorators)+len(decorators))
	next = append(next, f.decorators...)
	next = append(next, decorators...)
	return Factory[T]{decorators: next}
}

func (f Factory[T]) Wrap(cmd Command[T]) Command[T] {
	return Wrap(cmd, f.decorators...)
}

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
