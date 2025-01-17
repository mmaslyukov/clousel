package fault

type IError interface {
	Code() any
	Error() string
	Full() string
}
