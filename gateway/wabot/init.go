package wabot

import (
	"github.com/arpinfidel/tuduit/pkg/db"
)

type Server struct {
}

type Dependencies struct {
	DB *db.DB
}

func (s *Server) Start(deps Dependencies) (err error) {
	// https://godocs.io/go.mau.fi/whatsmeow#example-package
	return nil
}

func (s *Server) v1() {
}
