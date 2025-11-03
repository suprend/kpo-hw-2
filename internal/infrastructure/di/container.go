package di

import (
	"fmt"
	"reflect"
	"sync"
)

type Constructor func(c Container) (any, error)

type Container interface {
	Register(interfacePtr any, constructor Constructor) error
	Provide(interfacePtr any, instance any) error
	Resolve(interfacePtr any) (any, error)
}

type container struct {
	mu           sync.RWMutex
	constructors map[reflect.Type]Constructor
	singletons   map[reflect.Type]any
}

func New() Container {
	return &container{
		constructors: make(map[reflect.Type]Constructor),
		singletons:   make(map[reflect.Type]any),
	}
}

func (c *container) Register(interfacePtr any, constructor Constructor) error {
	t, err := extractType(interfacePtr)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.constructors[t] = constructor
	return nil
}

func (c *container) Provide(interfacePtr any, instance any) error {
	t, err := extractType(interfacePtr)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.singletons[t] = instance
	return nil
}

func (c *container) Resolve(interfacePtr any) (any, error) {
	t, err := extractType(interfacePtr)
	if err != nil {
		return nil, err
	}

	c.mu.RLock()
	instance, ok := c.singletons[t]
	if ok {
		c.mu.RUnlock()
		return instance, nil
	}

	constructor, ok := c.constructors[t]
	c.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("di: constructor for %s not registered", t)
	}

	instance, err = constructor(c)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if existing, exists := c.singletons[t]; exists {
		return existing, nil
	}

	c.singletons[t] = instance
	return instance, nil
}

func extractType(ptr any) (reflect.Type, error) {
	if ptr == nil {
		return nil, fmt.Errorf("di: nil type pointer")
	}

	t := reflect.TypeOf(ptr)
	if t.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("di: expected pointer to type, e.g. (*MyType)(nil)")
	}

	return t.Elem(), nil
}

func Register[T any](c Container, constructor func(Container) (T, error)) error {
	interfacePtr := (*T)(nil)
	return c.Register(interfacePtr, func(container Container) (any, error) {
		value, err := constructor(container)
		if err != nil {
			return nil, err
		}
		return value, nil
	})
}

func Provide[T any](c Container, instance T) error {
	interfacePtr := (*T)(nil)
	return c.Provide(interfacePtr, instance)
}

func Resolve[T any](c Container) (T, error) {
	interfacePtr := (*T)(nil)

	value, err := c.Resolve(interfacePtr)
	if err != nil {
		var zero T
		return zero, err
	}

	if value == nil {
		var zero T
		return zero, nil
	}

	typed, ok := value.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf(
			"di: resolved instance of type %T does not implement %s",
			value,
			reflect.TypeFor[T](),
		)
	}

	return typed, nil
}
