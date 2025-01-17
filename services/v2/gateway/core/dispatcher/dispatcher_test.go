package dispatcher_test

import (
	"bytes"
	"gateway/core/dispatcher"
	"testing"

	"github.com/google/uuid"
)

func TestUUIDMarshalCapabilities(t *testing.T) {
	var msg dispatcher.MessageGeneric
	msg.CarId = uuid.New()
	buf, _ := dispatcher.EncodeJson(&msg)
	t.Logf("buf: %s", buf.String())
	str := `{"Type":"New","CarId":"9552ffd4-fa96-4d12-89c3-4de416089bb3","SeqNum":10}`
	msgD, _ := dispatcher.DecodeJson[dispatcher.MessageGeneric](bytes.NewBuffer([]byte(str)))
	t.Logf("msg: %+v", msgD)

}
