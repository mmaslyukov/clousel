package machine_test

import (
	"clousel/core/machine"
	"testing"
)

type TestData struct {
	I   int    `json:"I"`
	Str string `json:"Str"`
}

func TestBinaryEncoderDecoder(t *testing.T) {
	a := TestData{I: 9999, Str: "123"}
	buf, err := machine.EncodeBinary(&a)
	if err != nil {
		t.Error(err.Error())
	}
	// t.Logf("a %+v", a)
	// t.Logf("%d -> %+v", len(buf.Bytes()), buf.Bytes())
	b, err := machine.DecodeBinary[TestData](&buf)
	if err != nil {
		t.Error(err.Error())
	}
	if a.I != b.I {
		t.Errorf("data mismatch")
	}
}
func TestJsonEncoderDecoder(t *testing.T) {
	a := TestData{I: 9999, Str: "123"}
	buf, err := machine.EncodeJson(&a)
	if err != nil {
		t.Error(err.Error())
	}
	b, err := machine.DecodeJson[TestData](&buf)
	if err != nil {
		t.Error(err.Error())
	}
	if a.I != b.I {
		t.Errorf("data mismatch")
	}
}
