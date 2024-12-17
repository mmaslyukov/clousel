package utils

import "fmt"

type Optional[T any] struct {
	value *T
}

func (o *Optional[T]) Set(value T) {
	*o.value = value
}

func (o *Optional[T]) Replace(value T) {
	o.value = &value
}

func (o *Optional[T]) Ptr() *T {
	return o.value
}

func (o *Optional[T]) Get() T {
	return *o.value
}

func (o *Optional[T]) IsValid() bool {
	return o.value != nil
}

func (o *Optional[T]) String() string {
	return fmt.Sprintf("%+v", o.value)
}

func OptionalValueCreate[T any](value T) Optional[T] {
	return Optional[T]{value: &value}
}

func OptionalNilCreate[T any]() Optional[T] {
	return Optional[T]{}
}
