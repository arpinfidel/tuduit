//go:generate sh -c "mockgen -source='../../../app/interface.go' -destination=$(basename $GOPACKAGE)_mock_test.go -package=$(basename $GOPACKAGE)"

package http

import (
	"net/http"

	taskrepo "github.com/arpinfidel/tuduit/repo/task"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Handler *Handler
}

func (s *Server) Start(deps Dependencies) (err error) {
	s.Handler = &Handler{
		deps: Dependencies{
			TaskRepo: *taskrepo.New(taskrepo.Dependencies{
				DB: deps.DB,
			}),
		},
	}
	r := chi.NewRouter()
	api := r.Route("/api", func(r chi.Router) {})
	api.Route("/v1", s.v1)

	return http.ListenAndServe(":2000", r)
}

func (s *Server) v1(r chi.Router) {
}
