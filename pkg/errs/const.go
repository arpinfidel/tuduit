package errs

import "errors"

var (
	ErrTypeBase     = errors.New("base")
	ErrTypeExpected = errors.New("expected")
)
