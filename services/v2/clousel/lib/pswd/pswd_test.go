package pswd_test

import (
	"clousel/lib/pswd"
	"testing"
)

func TestIMachineRepoGeneralAdapterInterface(t *testing.T) {
	p1 := pswd.PasswordPlainCreate("1234567890")
	t.Logf("p1.Plain %s", p1.Str())
	p2 := p1.Hash()
	t.Logf("p2.Hashed %s", p2.HexStr())
	// p2.Encode()
	p3 := p2.Encode()
	t.Logf("p3.Encoded %s", p3.HexStr())
	t.Logf("p3.EncodedStr %s", p3.Str())
	if p4, e := p3.Decode(); e == nil {
		t.Logf("p4.Decoded %s", p4.HexStr())
	} else {
		t.Error(e.Error())
	}

	t.Logf("Full line: %s", p1.Hash().Encode().Str())

}
