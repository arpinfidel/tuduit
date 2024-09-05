package errs

import (
	"errors"
	"net/http"
)

var (
	ErrBadRequest      = WrapHTTP(errors.New("bad request")).WithStatusCode(http.StatusBadRequest).WithAttributes(ErrTypeExpected)
	ErrUnauthorized    = WrapHTTP(errors.New("unauthorized")).WithStatusCode(http.StatusUnauthorized).WithAttributes(ErrBadRequest)
	ErrForbidden       = WrapHTTP(errors.New("forbidden")).WithStatusCode(http.StatusForbidden).WithAttributes(ErrBadRequest)
	ErrNotFound        = WrapHTTP(errors.New("not found")).WithStatusCode(http.StatusNotFound).WithAttributes(ErrBadRequest)
	ErrConflict        = WrapHTTP(errors.New("conflict")).WithStatusCode(http.StatusConflict).WithAttributes(ErrBadRequest)
	ErrTooManyRequests = WrapHTTP(errors.New("too many requests")).WithStatusCode(http.StatusTooManyRequests).WithAttributes(ErrBadRequest)
)

var (
	ErrInternalServerError     = WrapHTTP(errors.New("internal server error")).WithStatusCode(http.StatusInternalServerError).WithAttributes(ErrTypeBase)
	ErrServiceUnavailable      = WrapHTTP(errors.New("service unavailable")).WithStatusCode(http.StatusServiceUnavailable).WithAttributes(ErrInternalServerError)
	ErrGatewayTimeout          = WrapHTTP(errors.New("gateway timeout")).WithStatusCode(http.StatusGatewayTimeout).WithAttributes(ErrInternalServerError)
	ErrBadGateway              = WrapHTTP(errors.New("bad gateway")).WithStatusCode(http.StatusBadGateway).WithAttributes(ErrInternalServerError)
	ErrHTTPVersionNotSupported = WrapHTTP(errors.New("http version not supported")).WithStatusCode(http.StatusHTTPVersionNotSupported).WithAttributes(ErrInternalServerError)
	ErrVariantAlsoNegotiates   = WrapHTTP(errors.New("variant also negotiates")).WithStatusCode(http.StatusVariantAlsoNegotiates).WithAttributes(ErrInternalServerError)
	ErrInsufficientStorage     = WrapHTTP(errors.New("insufficient storage")).WithStatusCode(http.StatusInsufficientStorage).WithAttributes(ErrInternalServerError)
	ErrLoopDetected            = WrapHTTP(errors.New("loop detected")).WithStatusCode(http.StatusLoopDetected).WithAttributes(ErrInternalServerError)
	ErrNotExtended             = WrapHTTP(errors.New("not extended")).WithStatusCode(http.StatusNotExtended).WithAttributes(ErrInternalServerError)
)
