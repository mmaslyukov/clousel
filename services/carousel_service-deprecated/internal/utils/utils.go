package utils

import (
	"fmt"
	"sync"
)

type Optional[T any] struct {
	value *T
}

func (o *Optional[T]) Set(value T) {
	o.value = &value
}
func (o *Optional[T]) Ptr() *T {
	return o.value
}
func (o *Optional[T]) Get() T {
	return *o.value
}
func (o *Optional[T]) Valid() bool {
	return o.value != nil
}

func (o *Optional[T]) String() string {
	return fmt.Sprintf("%+v", o.value)
}
func NewOptionalValue[T any](value T) Optional[T] {
	return Optional[T]{value: &value}
}
func NewOptionalNil[T any]() Optional[T] {
	return Optional[T]{}
}

type Guard[T any] struct {
	mutex sync.RWMutex
	Value T
}

func (g *Guard[T]) Lock(f func()) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	f()
}
func (g *Guard[T]) RLock(f func()) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	f()
}
