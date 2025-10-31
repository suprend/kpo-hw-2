package di

import (
	"fmt"
	"reflect"
)

type Constructor func(c Container) (any, error)

type Container interface {
	Register(interfacePtr any, constructor Constructor)
	Provide(interfacePtr any, instance any)
	Resolve(interfacePtr any) (any, error)
}

type container struct {
	constructors map[reflect.Type]Constructor
	singletons   map[reflect.Type]any
}

func New() Container {
	return &container{
		constructors: make(map[reflect.Type]Constructor),
		singletons:   make(map[reflect.Type]any),
	}
}

func (c *container) Register(interfacePtr any, constructor Constructor) {
	t := mustExtractInterface(interfacePtr)
	c.constructors[t] = constructor
}

func (c *container) Provide(interfacePtr any, instance any) {
	t := mustExtractInterface(interfacePtr)
	c.singletons[t] = instance
}

func (c *container) Resolve(interfacePtr any) (any, error) {
	t, err := extractInterface(interfacePtr)
	if err != nil {
		return nil, err
	}

	if instance, ok := c.singletons[t]; ok {
		return instance, nil
	}

	constructor, ok := c.constructors[t]
	if !ok {
		return nil, fmt.Errorf("di: constructor for %s not registered", t)
	}

	instance, err := constructor(c)
	if err != nil {
		return nil, err
	}

	c.singletons[t] = instance
	return instance, nil
}

func extractInterface(ptr any) (reflect.Type, error) {
	if ptr == nil {
		return nil, fmt.Errorf("di: nil interface pointer")
	}

	t := reflect.TypeOf(ptr)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Interface {
		return nil, fmt.Errorf("di: expected pointer to interface, e.g. (*SomeInterface)(nil)")
	}

	return t.Elem(), nil
}

func mustExtractInterface(ptr any) reflect.Type {
	t, err := extractInterface(ptr)
	if err != nil {
		panic(err)
	}
	return t
}
