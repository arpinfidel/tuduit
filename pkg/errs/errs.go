package errs

import (
	"errors"
	"fmt"
	"runtime"
)

var _ error = &Error{}

type Error struct {
	Base       error
	Attributes []error
	Trace      []string
}

func (e *Error) Error() string {
	return e.Base.Error()
}

func New(format string, a ...any) *Error {
	return &Error{
		Base:       fmt.Errorf(format, a...),
		Attributes: []error{ErrTypeBase},
	}
}

func Wrap(err error) *Error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrTypeBase) {
		return err.(*Error)
	}
	return &Error{
		Base:       err,
		Attributes: []error{ErrTypeBase},
	}
}

func TraceSkip(err error, skip int) *Error {
	if err == nil {
		return nil
	}
	return Wrap(err).WithTraceSkip(1 + skip)
}

func Trace(err error) *Error {
	if err == nil {
		return nil
	}
	return Wrap(err).WithTraceSkip(1)
}

func GetTrace(err error) []string {
	if !errors.Is(err, ErrTypeBase) {
		return nil
	}

	return err.(*Error).Trace
}

func (e *Error) WithTraceSkip(skip int) *Error {
	e.Trace = createStackTrace(1 + skip)
	return e
}

func (e *Error) WithTrace() *Error {
	return e.WithTraceSkip(1)
}

func (e *Error) WithAttributes(a ...error) *Error {
	e.Attributes = append(e.Attributes, a...)
	return e
}

func (e *Error) Unwrap() []error {
	return append([]error{e.Base}, e.Attributes...)
}

func createStackTrace(ignore int) []string {
	trace := []string{}
	pc := make([]uintptr, 15)
	n := runtime.Callers(0, pc)
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, next := frames.Next()
		t := fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function)
		trace = append(trace, t)
		if !next {
			break
		}
	}
	return trace[2+ignore:]
}

func DeferTrace(err *error) func() error {
	return func() error {
		if err == nil || *err == nil {
			return nil
		}

		*err = TraceSkip(*err, 1)
		return *err
	}
}
