package trace

import (
	"context"

	"github.com/arpinfidel/tuduit/pkg/errs"
)

func Default(ctxAddr *context.Context, errAddr *error) func() {
	return func() {
		errs.WrappError(errAddr)
	}
}
