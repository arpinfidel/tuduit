package errs

import (
	"errors"
	"fmt"
	"runtime"
)

var _ error = &Err{}

type Err struct {
	Base       error
	Attributes []error
	Trace      []string

	UserMessage string
}

func (e *Err) Error() string {
	return e.Base.Error()
}

func New(format string, a ...any) *Err {
	return &Err{
		Base:       fmt.Errorf(format, a...),
		Attributes: []error{ErrTypeBase},
	}
}

func Wrap(err error) *Err {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrTypeBase) {
		return err.(*Err)
	}
	return &Err{
		Base:       err,
		Attributes: []error{ErrTypeBase},
	}
}

func TraceSkip(err error, skip int) *Err {
	if err == nil {
		return nil
	}
	return Wrap(err).WithTraceSkip(1 + skip)
}

func Trace(err error) *Err {
	if err == nil {
		return nil
	}
	return Wrap(err).WithTraceSkip(1)
}

func GetTrace(err error) []string {
	if !errors.Is(err, ErrTypeBase) {
		return nil
	}

	return err.(*Err).Trace
}

func (e *Err) WithTraceSkip(skip int) *Err {
	e.Trace = createStackTrace(1 + skip)
	return e
}

func (e *Err) WithTrace() *Err {
	return e.WithTraceSkip(1)
}

func (e *Err) WithAttributes(a ...error) *Err {
	e.Attributes = append(e.Attributes, a...)
	return e
}

func (e *Err) WithUserMessagef(format string, a ...any) *Err {
	e.UserMessage = fmt.Sprintf(format, a...)
	return e
}

func (e *Err) Unwrap() []error {
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
