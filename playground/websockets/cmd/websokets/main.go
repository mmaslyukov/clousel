package main

/*
Controller <>---> HadlerPortInterface
                       .
                      /_\
                       |
                   HandlerPort
                       |
                       V
                     Logic <>----> AdapterPortInterface
                                           .
                     				      /_\
                     				       |
                     			       AdapterPort
*/

import (
	// "fmt"
	"log"
	"net/http"
	"playground/wse"
	"sync"
)

func NewGuard[T any](resourse T) Guard[T] {
	return Guard[T]{resource: resourse}
}

type Guard[T any] struct {
	sync.RWMutex
	resource T
}

func (g *Guard[T]) LockGet() *T {
	g.Lock()
	return &g.resource
}

// func (g *Guard[T]) Unlock() {
// 	g.Unlock()
// }

// func (g *Guard[T]) Open() Lock[T] {
// 	return Lock[T]{mutex: g, resource: &g.resource}
// }

// type Lock[T any] struct {
// 	mutex    *sync.Mutex
// 	resource *T
// }

//	func (l *Lock[T]) Get() *T {
//		return l.resource
//	}
//
//	func (l *Lock[T]) Close() {
//		l.mutex.Unlock()
//	}

type HandlerPortReceiverInterface[T any] interface {
	Receiver() chan T
}
type HandlerPortTransmitterInterface[T any] interface {
	Transmitter() chan T
}

type HandlerPortTransmitter[T any] struct {
	tx chan T
}

func (h *HandlerPortTransmitter[T]) Transmitter() chan T {
	return h.tx
}

var g Guard[int]

// func main() {
// 	g = NewGuard[int](5)
// 	go func() {
// 		for {
// 			time.Sleep(100 * time.Millisecond)
// 			*g.LockGet() += 1
// 			g.Unlock()
// 		}
// 	}()
// 	go func() {
// 		for {
// 			time.Sleep(500 * time.Millisecond)
// 			*g.LockGet() += 100
// 			g.Unlock()
// 		}
// 	}()
// 	for {
// 		time.Sleep(100 * time.Millisecond)
// 		fmt.Println(*g.LockGet())
// 		g.Unlock()
// 	}
// 	// wse.Foo()
// 	time.Sleep(10 * time.Second)
// 	fmt.Println("Done")
// }

func main() {
	core := wse.NewCore()
	go core.Run()
	// http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wse.ServeWs(core, w, r)
	})

	log.Print("Starting server...")

	log.Fatal(http.ListenAndServe("localhost:8080", nil))

}
