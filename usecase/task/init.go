package taskuc

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
		deps: deps,
	}
}

type IRepo interface {
	repo.IStdRepo[entity.Task]
	IExtRepo
}

type IExtRepo interface {
}
