package ports

type PortInterface[T any] interface {
	Send(value T)
	Receiver() <-chan T
}

type Port[T any] struct {
	c chan T
}

func (p *Port[T]) Send(value T) {
	p.c <- value
}
func (p *Port[T]) Receiver() <-chan T {
	return p.c
}

func NewPort[T any](size int) *Port[T] {
	return &Port[T]{c: make(chan T, size)}
}

// type PortReceiverInterface[T any] interface {
// 	Receiver() chan T
// }
// type PortTransmitterInterface[T any] interface {
// 	Transmitter() chan T
// }

// type PortReceiverHandler[T any] struct {
// 	rx chan T
// }

// func NewPortReceiverHandler[T any](size int) *PortReceiverHandler[T] {
// 	return &PortReceiverHandler[T]{rx: make(chan T, size)}
// }

// func (h *PortReceiverHandler[T]) Receiver() chan T {
// 	return h.rx
// }

// type PortTransmitterAdapter[T any] struct {
// 	tx chan T
// }

// func NewPortTransmitterAdapter[T any](size int) *PortTransmitterAdapter[T] {
// 	return &PortTransmitterAdapter[T]{tx: make(chan T)}
// }
// func (h *PortTransmitterAdapter[T]) Transmitter() chan T {
// 	return h.tx
// }
