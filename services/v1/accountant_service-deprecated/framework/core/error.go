package core

import "fmt"

type Error[T any] struct {
	message string
}

func (e *Error[T]) Error() string {
	return e.message
}

func (e *Error[T]) Message(format string, a ...any) *Error[T] {
	if len(a) > 0 {
		e.message = fmt.Sprintf("%s: %s", e.message, fmt.Sprintf(format, a...))
	} else {
		e.message = fmt.Sprintf("%s: %s", e.message, fmt.Sprint(format))
	}
	return e
}

func ErrorCreate[T any]() *Error[T] {
	return &Error[T]{message: fmt.Sprintf("%T", *new(T))}
}
