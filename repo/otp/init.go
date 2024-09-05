package otprepo

import (
	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/arpinfidel/tuduit/pkg/log"
	"github.com/arpinfidel/tuduit/repo"
)

type Repo struct {
	deps Dependencies

	repo.DBConnection
	*repo.StdCRUD[entity.OTP]
}

type Dependencies struct {
	DB     *db.DB
	Logger *log.Logger
}

func New(deps Dependencies) *Repo {
	return &Repo{
		deps:         deps,
		DBConnection: *repo.NewDBConnection(deps.DB),
		StdCRUD:      repo.NewStdCRUD[entity.OTP](deps.DB, deps.Logger, "trx_otp"),
	}
}
