package entity

import (
	"database/sql/driver"

	"github.com/lib/pq"
)

type PQArr[T any] []T

func (a *PQArr[T]) Scan(src interface{}) error {
	pqArr := pq.Array((*[]T)(a))
	return pqArr.Scan(src)
}

func (a PQArr[T]) Value() (val driver.Value, err error) {
	pqArr := pq.Array(([]T)(a))
	return pqArr.Value()
}
