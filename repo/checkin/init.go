package checkinrepo

import (
	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/arpinfidel/tuduit/pkg/log"
	"github.com/arpinfidel/tuduit/repo"
)

type Repo struct {
	deps Dependencies

	repo.DBConnection
	*repo.StdCRUD[entity.CheckIn]
}

type Dependencies struct {
	DB     *db.DB
	Logger *log.Logger
}

func New(deps Dependencies) *Repo {
	return &Repo{
		deps:         deps,
		DBConnection: *repo.NewDBConnection(deps.DB),
		StdCRUD:      repo.NewStdCRUD[entity.CheckIn](deps.DB, deps.Logger, "mst_check_in"),
	}
}
