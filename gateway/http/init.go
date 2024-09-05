//go:generate sh -c "mockgen -source='../../../app/interface.go' -destination=$(basename $GOPACKAGE)_mock_test.go -package=$(basename $GOPACKAGE)"

package http

import (
	"errors"
	"net/http"

	"github.com/arpinfidel/tuduit/app"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/arpinfidel/tuduit/pkg/errs"
	"github.com/arpinfidel/tuduit/pkg/log"
	"github.com/arpinfidel/tuduit/pkg/rose"
	"github.com/go-chi/chi/v5"
)

type HTTPServer struct {
	l *log.Logger

	d Dependencies
}

type Dependencies struct {
	App *app.App
}

func New(l *log.Logger, deps Dependencies) *HTTPServer {
	return &HTTPServer{
		l: l,
		d: deps,
	}
}

func (s *HTTPServer) Start() (err error) {
	r := chi.NewRouter()
	s.route(r)

	return http.ListenAndServe(":2000", r)
}

type PublicError struct {
	UserMessage string            `rose:"user_message"`
	Fields      map[string]string `rose:"fields"`
}

type PrivateError struct {
	DevMessage string   `rose:"dev_message"`
	Trace      []string `rose:"trace"`
}

type Error struct {
	PublicError
	PrivateError
}

type JSONResponse[T any] struct {
	Data T      `rose:"data"`
	Err  *Error `rose:"err,omitempty"`

	status int
}

func NewJSONResponse[T any](data T) *JSONResponse[T] {
	return &JSONResponse[T]{
		Data: data,

		status: http.StatusOK,
	}
}

func NewJSONResponseErr[T any](data T, err error) *JSONResponse[T] {
	resp := JSONResponse[T]{
		Data: data,
		Err:  &Error{},

		status: http.StatusInternalServerError,
	}

	resp.Err.DevMessage = err.Error()
	var e *errs.Err
	if errors.As(err, &e) {
		resp.Err.UserMessage = e.UserMessage
		resp.Err.Trace = e.Trace
	}

	var eh *errs.HTTPErrror
	if errors.As(err, &eh) {
		resp.status = eh.StatusCode
	}

	return &resp
}

func (j *JSONResponse[T]) Response(w http.ResponseWriter, r *http.Request) *JSONResponse[T] {
	w.WriteHeader(j.status)
	w.Header().Set("Content-Type", "application/json")
	b, err := rose.JSONMarshal(j)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("INTERNAL SERVER ERROR"))
		return j
	}

	w.Write(b)
	return j
}

func (j *JSONResponse[T]) Log() *JSONResponse[T] {
	// TODO:

	return j
}

func wrapHandler[T, U any](s *HTTPServer, f func(ctx *ctxx.Context, t T) (U, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := ctxx.GetContext(r.Context())

		var t T
		rose, err := rose.NewParser(ctx, "").ParseHTTP(r, &t)
		if err != nil {
			NewJSONResponseErr(t, errs.ErrBadRequest.WithTrace().WithUserMessagef("invalid request: %v", err)).Response(w, r)
			return
		}
		if !rose.Valid {
			NewJSONResponseErr(t, errs.ErrBadRequest.WithTrace().WithUserMessagef("invalid request: %v", rose.Errors[0])).Response(w, r) // TODO: handle multiple errors
			return
		}

		resp, err := f(ctx, t)
		if err != nil {
			s.l.Errorf("failed to process request: %v", err)
			NewJSONResponseErr(t, err).Response(w, r)
			return
		}

		NewJSONResponse(resp).Response(w, r)
	}
}
