package errs

import (
	"net/http"
)

type HTTPErrror struct {
	*Err

	StatusCode int
}

func (e *HTTPErrror) WithStatusCode(statusCode int) *HTTPErrror {
	e.StatusCode = statusCode
	return e
}

func WrapHTTP(err error) *HTTPErrror {
	if err == nil {
		return nil
	}

	return &HTTPErrror{
		Err:        Wrap(err),
		StatusCode: http.StatusInternalServerError,
	}
}
