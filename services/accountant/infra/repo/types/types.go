package types

type Named[T any] struct {
	name  string
	Value T
}

func (n *Named[T]) Name() string {
	return n.name
}

func NamedCreateDefault[T any](name string) Named[T] {
	return Named[T]{name: name}
}

type NamedOpt[T any] struct {
	name  string
	Value *T
}

func (n *NamedOpt[T]) Name() string {
	return n.name
}

func NamedOptCreateDefault[T any](name string) NamedOpt[T] {
	return NamedOpt[T]{name: name}
}
