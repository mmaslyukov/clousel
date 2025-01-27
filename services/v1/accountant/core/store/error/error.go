package error

import "fmt"

type ErrorCode int

const (
	ECGeneralFailure ErrorCode = iota
	ECReadKeys
	ECPayment
	ECBookRepoRead
	ECBookRepoInsert
	ECBookRepoMark
	ECRemoteServiceCarouselRefill
	ECStripeCheckout
	ECStripeReadPrice
	ECStripeReadProduct
	ECStripeCheckoutSession
)

type IError interface {
	Code() ErrorCode
	Full() string
	Error() string
}

type Error struct {
	code ErrorCode
	msg  string
	err  IError
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Code() ErrorCode {
	return e.code
}
func (e *Error) Err(err IError) *Error {
	e.err = err
	return e
}
func (e *Error) Msgf(format string, v ...interface{}) *Error {
	e.msg = fmt.Sprintf(format, v...)
	return e
}
func (e *Error) Msg(msg string) *Error {
	e.msg = msg
	return e
}
func (e *Error) Full() string {
	imsg := ""
	if e.err != nil {
		imsg = e.err.Full()
	}
	msg := fmt.Sprintf(">> EC(%d) %s %s ", e.code, e.msg, imsg)

	return msg
}

func New(code ErrorCode) *Error {
	return &Error{code: code}
}
