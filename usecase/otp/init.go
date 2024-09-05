package otpuc

import (
	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/repo"
)

type UseCase struct {
	IRepo
	deps Dependencies
}

type Dependencies struct {
	Repo IRepo
}

func New(deps Dependencies) *UseCase {
	return &UseCase{
		IRepo: deps.Repo,
		deps:  deps,
	}
}

type IRepo interface {
	repo.IStdRepo[entity.OTP]
	IExtRepo
}

type IExtRepo interface {
}
