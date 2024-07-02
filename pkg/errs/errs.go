package errs

import (
	"errors"
	"fmt"
	"runtime"
)

// Error struct
type Error struct {
	message string
	cause   error
	trace   []string
}

func Cause(err error) error {
	if err, ok := err.(Error); ok {
		return err.cause
	}
	return err
}

func createStackTrace(skip int) []string {
	trace := []string{}
	pc := make([]uintptr, 15)
	n := runtime.Callers(skip+2, pc)
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, next := frames.Next()
		t := fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function)
		trace = append(trace, t)
		if !next {
			break
		}
	}
	return trace
}

// NewError creates a new error and returns itself for any additional functions like setting the message
func NewError(str string) *Error {
	return &Error{
		cause: errors.New(str),
		trace: createStackTrace(0),
	}
}

// WrapError wraps the original error
func WrapError(err error) *Error {
	return &Error{
		cause: err,
		trace: createStackTrace(0),
	}
}

// WrappError wraps the error in place
func WrappError(err *error) *Error {
	if err == nil || *err == nil {
		return nil
	}

	if e, ok := (*err).(*Error); ok {
		return e
	}

	e := &Error{
		cause: *err,
		trace: createStackTrace(0),
	}

	*err = e

	return e
}

// WrapfError wraps the original error and sets a custom message
func WrapfError(err error, format string, v ...interface{}) *Error {
	txt := fmt.Sprintf(format, v...)
	return &Error{
		cause:   err,
		message: txt,
		trace:   createStackTrace(0),
	}
}

// Error implements the error interface
func (e Error) Error() string {
	msg := e.cause.Error()
	if e.message != "" {
		msg = fmt.Sprintf("%s: %s", e.message, msg)
	}
	return msg
}

// Cause return original error
func (e Error) Cause() error {
	return e.cause
}

// Trace returns the stack trace
func (e *Error) Trace() []string {
	return e.trace
}
