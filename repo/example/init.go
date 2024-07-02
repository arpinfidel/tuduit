package example

import (
	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/repo"
)

type Repo struct {
	deps Dependencies
	repo.StdCRUD[entity.Example]
}

type Dependencies struct {
	DB repo.DBConnection
}

func New(deps Dependencies) *Repo {
	return &Repo{
		deps: deps,
	}
}
