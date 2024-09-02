package userrepo

import (
	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/arpinfidel/tuduit/pkg/log"
	"github.com/arpinfidel/tuduit/repo"
)

type Repo struct {
	deps Dependencies

	repo.DBConnection
	*repo.StdCRUD[entity.User]
}

type Dependencies struct {
	DB     *db.DB
	Logger *log.Logger
}

func New(deps Dependencies) *Repo {
	return &Repo{
		deps:         deps,
		DBConnection: *repo.NewDBConnection(deps.DB),
		StdCRUD:      repo.NewStdCRUD[entity.User](deps.DB, deps.Logger, "mst_user"),
	}
}
