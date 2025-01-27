package main

import (
	"accountant_service/framework/core"
	"fmt"
)

type Foo struct {
}

func main() {
	err := core.ErrorCreate[Foo]().Message("Test")
	fmt.Println(err.Error())
	err = core.ErrorCreate[Foo]().Message("Test %d, '%s'", 1, "fuck")
	// fmt.Printf("\nTest %d, '%s'\n", 1, "fuck")

	fmt.Println(err.Error())
	// fmt.Println(reflect.TypeOf(Foo).String())
	// i := internal.DummyTwo{}
}
