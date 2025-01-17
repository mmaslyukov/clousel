package fault_test

import (
	"fmt"
	"gateway/lib/fault"
	"testing"
)

const (
	ErrorOne = iota
	ErrorTwo
)

func TestError(t *testing.T) {
	e2 := fmt.Errorf("My general error")
	e1 := fault.New(ErrorOne).Err(e2).Msg("My Error One fault")
	t.Logf("Err Text is: %s", e1.Full())
}
