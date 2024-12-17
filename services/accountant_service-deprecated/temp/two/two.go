package two

import (
	"accountant_service/temp/one"
	"fmt"
)

type Two struct {
	one.One
}

func (o *Two) PrintTwo() {
	fmt.Println("Two")
}
