package error_test

import (
	"accountant/core/owner/error"
	"testing"
)

func TestError(t *testing.T) {
	e1 := error.New(error.ECCarouselRegisterFailure).Msgf("TestError: ECCarouselRegisterFailure")
	e2 := error.New(error.ECProductAssignFailure).Msgf("TestError: ECProductAssignFailure").Err(e1)
	t.Logf("Err Text is: %s", e2.Full())
	// fmt.Printf("Err Text is: %s", e2.Full())
	// t.Errorf("Expected '%s', but got '%s'", expect, tp.Get())
}
