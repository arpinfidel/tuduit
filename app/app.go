package app

import "context"

type Context struct {
	context.Context

	UserID int64
}
