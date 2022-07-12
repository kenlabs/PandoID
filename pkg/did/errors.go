package did

import "fmt"

var ErrInvalidDID = ParserError{msg: "invalid DID"}

type ParserError struct {
	msg string
	err error
}

func (w ParserError) wrap(err error) error {
	return ParserError{msg: fmt.Sprintf("%s: %s", w.msg, err.Error()), err: err}
}

func (w ParserError) Is(other error) bool {
	_, ok := other.(ParserError)
	return ok
}

func (w ParserError) Unwrap() error {
	return w.err
}

func (w ParserError) Error() string {
	return w.msg
}
