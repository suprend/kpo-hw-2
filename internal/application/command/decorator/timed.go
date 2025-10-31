package decorator

import (
	"context"
	"log"
	"time"

	"kpo-hw-2/internal/application/command"
)

// Timed decorates commands, measuring execution time and logging it.
type Timed[T any] struct {
	Log   func(name string, duration time.Duration, err error)
	Clock func() time.Time
}

// Wrap implements command.Decorator.
func (d Timed[T]) Wrap(inner command.Command[T]) command.Command[T] {
	if inner == nil {
		return nil
	}

	logFn := d.Log
	if logFn == nil {
		logFn = func(name string, duration time.Duration, err error) {
			log.Printf("%s took %s (err=%v)", name, duration, err)
		}
	}

	clock := d.Clock
	if clock == nil {
		clock = time.Now
	}
	name := inner.Name()

	return command.Func[T]{
		NameFn: func() string { return name },
		ExecFn: func(ctx context.Context) (T, error) {
			start := clock()
			result, err := inner.Execute(ctx)
			finish := clock()

			logFn(name, finish.Sub(start), err)
			return result, err
		},
	}
}
