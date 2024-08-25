package taskrepo

import (
	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/arpinfidel/tuduit/repo"
)

type Repo struct {
	deps Dependencies

	repo.DBConnection
	*repo.StdCRUD[entity.Task]
}

type Dependencies struct {
	DB *db.DB
}

func New(deps Dependencies) *Repo {
	return &Repo{
		deps:         deps,
		DBConnection: *repo.NewDBConnection(deps.DB),
		StdCRUD:      repo.NewStdCRUD[entity.Task](deps.DB, "mst_task"),
	}
}
