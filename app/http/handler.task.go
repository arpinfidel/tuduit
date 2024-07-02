package http

import (
	"context"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/arpinfidel/tuduit/pkg/errs"
	"github.com/arpinfidel/tuduit/pkg/trace"
)

func (h Handler) getTasks(
	ctx context.Context,
	path struct{},
	query struct {
		entity.Pagination
	},
	req struct{},
) (book entity.Task, err error) {
	defer trace.Default(&ctx, &err)

	books, _, err := h.deps.TaskRepo.Get(ctx, nil, db.Params{
		Pagination: query.QBPaginate(),
	})
	if len(books) == 0 {
		return entity.Task{}, errs.NewError("task not found")
	}
	return books[0], err
}

func (h Handler) getTaskByID(
	ctx context.Context,
	path struct {
		ID string `json:"id"`
	},
	query struct{},
	req struct{},
) (book entity.Task, err error) {
	defer trace.Default(&ctx, &err)

	books, _, err := h.deps.TaskRepo.GetByIDs(ctx, nil, []string{path.ID}, entity.Pagination{Page: 1, PageSize: 1})
	if len(books) == 0 {
		return entity.Task{}, errs.NewError("task not found")
	}
	return books[0], err
}

func (h Handler) newTask(
	ctx context.Context,
	path struct{},
	query struct{},
	req []entity.Task,
) (books []entity.Task, err error) {
	defer trace.Default(&ctx, &err)

	books, err = h.deps.TaskRepo.Create(ctx, nil, req)
	return books, err
}
