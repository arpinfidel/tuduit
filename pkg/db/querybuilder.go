package db

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

type Params struct {
	Where      []Where
	Pagination *Pagination
	Sort       []Sort

	where string
	page  string
	sort  string

	args []any
}

type Where struct {
	Field  string
	Op     Operator
	Value  any
	RawSQL string
}

type Operator int

const (
	DefaultOp Operator = iota
	NotNullOp
	IsNullOp
	EqOp
	NotEqOp
	GtOp
	LtOp
	GtOrEqOp
	LtOrEqOp
	InOp
	LikeOp
	OrOp
	AndOp
	RawOp
)

func (w *Where) buildOperator() sq.Sqlizer {
	v := w.Value
	switch w.Op {
	default:
		return sq.Eq{w.Field: v}
	case NotNullOp:
		return sq.NotEq{w.Field: nil}
	case IsNullOp:
		return sq.Eq{w.Field: v}
		// switch vv := v.(type) {
		// default:
		// 	panic(fmt.Sprintf("unsupported type %T", v))
		// case []string:
		// 	return in(w.Field, vv)
		// case []int:
		// 	return in(w.Field, vv)
		// case []int64:
		// 	return in(w.Field, vv)
		// case nil:
		// 	return sq.Expr("1 = 0")
		// }
	case EqOp:
		return sq.Eq{w.Field: v}
	case NotEqOp:
		return sq.NotEq{w.Field: v}
	case GtOp:
		return sq.Gt{w.Field: v}
	case LtOp:
		return sq.Lt{w.Field: v}
	case GtOrEqOp:
		return sq.GtOrEq{w.Field: v}
	case LtOrEqOp:
		return sq.LtOrEq{w.Field: v}
	case InOp:
		return sq.Eq{w.Field: v}
	case LikeOp:
		return sq.Like{w.Field: v}
	case OrOp:
		return or(v)
	case AndOp:
		return and(v)
	case RawOp:
		return sq.Expr(w.RawSQL)
	}
}

func in[T any](field string, arr []T) sq.Or {
	in := sq.Or{}
	for _, v := range arr {
		in = append(in, sq.Eq{field: v})
	}
	return in
}

func or(v interface{}) sq.Or {
	orFilters := sq.Or{}
	switch filters := v.(type) {
	case []Where:
		for _, f := range filters {
			orFilters = append(orFilters, f.buildOperator())
		}
	}

	return orFilters
}

func and(v interface{}) sq.And {
	orFilters := sq.And{}
	switch filters := v.(type) {
	case []Where:
		for _, f := range filters {
			orFilters = append(orFilters, f.buildOperator())
		}
	}

	return orFilters
}

func (p *Params) BuildWhere() *Params {
	if len(p.Where) == 0 {
		return p
	}

	filters := sq.And{}
	for _, w := range p.Where {
		filters = append(filters, w.buildOperator())
	}

	p.where = ""
	q, args, err := sq.Select("*").From("x").Where(filters).ToSql()
	if err != nil {
		panic(err)
	}

	// trim "SELECT * FROM x ..."
	p.where = q[16:]
	p.args = append(p.args, args...)

	return p
}

type Pagination struct {
	Limit  int
	Offset int
}

func (p *Params) BuildPagination() *Params {
	if p.Pagination == nil {
		return p
	}

	q, args, err := sq.Select("*").From("x").Limit(uint64(p.Pagination.Limit)).Offset(uint64(p.Pagination.Offset)).ToSql()
	if err != nil {
		panic(err)
	}

	p.page = q[16:]
	p.args = append(p.args, args...)

	return p
}

type Sort struct {
	Field      string
	Asc        bool
	NullsFirst bool
}

func (p *Params) BuildSort() *Params {
	if len(p.Sort) == 0 {
		return p
	}

	orderBy := []string{}
	for _, sort := range p.Sort {
		dir := "ASC"
		if !sort.Asc {
			dir = "DESC"
		}

		nulls := "FIRST"
		if !sort.NullsFirst {
			nulls = "LAST"
		}
		orderBy = append(orderBy, pq.QuoteIdentifier(sort.Field)+" "+dir+" NULLS "+nulls)
	}

	q, args, err := sq.Select("*").From("x").OrderBy(orderBy...).ToSql()
	if err != nil {
		panic(err)
	}

	// trim "SELECT * FROM x ..."
	p.sort = q[16:]
	p.args = append(p.args, args...)

	return p
}

func (p *Params) GetQuery(q string) (string, []interface{}) {
	return fmt.Sprintf("%s %s %s %s", q, p.where, p.sort, p.page), p.args
}

func (p *Params) GetCountQuery(q string) (string, []interface{}) {
	return fmt.Sprintf(q, p.where), p.args
}
