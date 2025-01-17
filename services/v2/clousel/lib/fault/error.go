package fault

import "fmt"

type Error struct {
	code any
	msg  string
	err  IError
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Code() any {
	return e.code
}
func (e *Error) Err(err error) *Error {
	msg := fmt.Sprint(err.Error())
	e.msg = e.concat(e.msg, msg)
	return e
}
func (e *Error) Fault(err IError) *Error {
	e.err = err
	return e
}
func (e *Error) Msgf(format string, v ...interface{}) *Error {
	msg := fmt.Sprintf(format, v...)
	e.msg = e.concat(e.msg, msg)

	return e
}
func (e *Error) Msg(msg string) *Error {
	e.msg = e.concat(e.msg, msg)
	return e
}
func (e *Error) Full() string {
	imsg := ""
	if e.err != nil {
		imsg = e.err.Full()
	}
	msg := fmt.Sprintf("\n>> EC:%v %s %s ", e.code, e.msg, imsg)
	return msg
}
func (e *Error) concat(a, b string) string {
	if len(a) > 0 && len(b) > 0 {
		return fmt.Sprintf("%s; %s", a, b)
	} else {
		return fmt.Sprintf("%s%s", a, b)
	}
}

func New(code any) *Error {
	return &Error{code: code}
}
