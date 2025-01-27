package operator

type OperatorError struct {
	str    string
	status int
}

func (e *OperatorError) Error() string {
	return e.str
}

func (e *OperatorError) Status() int {
	return e.status
}
