package repo

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/arpinfidel/tuduit/pkg/trace"
	"github.com/jmoiron/sqlx"
)

type IStdRepo[T any] interface {
	Create(ctx context.Context, dbTx *sqlx.Tx, newData []T) (data []T, err error)
	Update(ctx context.Context, dbTx *sqlx.Tx, newData T) (data T, err error)
	Delete(ctx context.Context, dbTx *sqlx.Tx, id string) (err error)
	Get(ctx context.Context, dbTx *sqlx.Tx, param db.Params) (data []T, total int, err error)
	GetByIDs(ctx context.Context, dbTx *sqlx.Tx, ids []int, pg entity.Pagination) (data []T, total int, err error)
}

type StdFields interface {
	GetStdFields() entity.StdFields
}

type StdCRUD[T StdFields] struct {
	db *db.DB

	tableName string

	fields             []string
	fieldsString       string
	placeholders       []string
	placeholdersString string
	setString          string
}

func NewStdCRUD[T StdFields](db *db.DB, tableName string) *StdCRUD[T] {
	return &StdCRUD[T]{
		db:        db,
		tableName: tableName,
	}
}

func (c *StdCRUD[T]) getFields() []string {
	if c.fields != nil {
		return c.fields
	}
	var x T
	elem := reflect.TypeOf(x)
	kind := elem.Kind()
	if kind == reflect.Array ||
		kind == reflect.Chan ||
		kind == reflect.Map ||
		kind == reflect.Ptr ||
		kind == reflect.Slice {
		elem = elem.Elem()
	}
	fieldN := elem.NumField()
	for i := 0; i < fieldN; i++ {
		field := elem.Field(i)
		tag := string(field.Tag.Get("db"))
		if field.Type == reflect.TypeOf(entity.StdFields{}) {

			c.fields = append(c.fields, c.getStdFields()...)
		} else {
			c.fields = append(c.fields, tag)
		}
	}

	return c.fields
}

var stdFields []string = nil

func (c *StdCRUD[T]) getStdFields() []string {
	if stdFields != nil {
		return stdFields
	}

	t := reflect.TypeOf(entity.StdFields{})
	var fields []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")
		if tag == "id" {
			continue
		}
		fields = append(fields, tag)
	}

	stdFields = fields

	return fields
}

func (c *StdCRUD[T]) getFieldsString() string {
	if c.fieldsString != "" {
		return c.fieldsString
	}
	c.fieldsString = strings.Join(c.getFields(), ",")
	return c.fieldsString
}

func (c *StdCRUD[T]) getPlaceholders() []string {
	if c.placeholders != nil {
		return c.placeholders
	}
	for _, f := range c.getFields() {
		c.placeholders = append(c.placeholders, ":"+f)
	}
	return c.placeholders
}

func (c *StdCRUD[T]) getPlaceholdersString() string {
	if c.placeholdersString != "" {
		return c.placeholdersString
	}
	c.placeholdersString = strings.Join(c.getPlaceholders(), ",")
	return c.placeholdersString
}

func (c *StdCRUD[T]) getSetString() string {
	if c.setString != "" {
		return c.setString
	}
	var set []string
	for _, f := range c.fields {
		set = append(set, fmt.Sprintf("%s=:%s", f, f))
	}
	c.setString = strings.Join(set, ",")
	return c.setString
}

func (c *StdCRUD[T]) Create(ctx context.Context, dbTx *sqlx.Tx, newData []T) (data []T, err error) {
	defer trace.Default(&ctx, &err)()

	var querier db.Querier = dbTx
	if dbTx == nil {
		querier = c.db.GetMaster()
	}

	fields := c.getFieldsString()
	placeHolders := c.getPlaceholdersString()

	q := `INSERT INTO %s (%s) VALUES (%s) RETURNING id,%s`
	q = fmt.Sprintf(q, c.tableName, fields, placeHolders, fields)

	q, arg, err := sqlx.Named(q, newData)
	if err != nil {
		return newData, err
	}

	err = querier.SelectContext(ctx, &data, c.db.Rebind(q), arg...)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (c *StdCRUD[T]) Update(ctx context.Context, dbTx *sqlx.Tx, newData T) (data T, err error) {
	defer trace.Default(&ctx, &err)()

	var querier db.Querier = dbTx
	if dbTx == nil {
		querier = c.db.GetMaster()
	}

	q := `UPDATE %s SET %s WHERE id=:id RETURNING id,%s`
	q = fmt.Sprintf(q, c.tableName, c.getSetString(), c.getFieldsString())

	q, arg, err := sqlx.Named(q, newData)
	if err != nil {
		return newData, err
	}

	err = querier.GetContext(ctx, &data, c.db.Rebind(q), arg...)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (c *StdCRUD[T]) Delete(ctx context.Context, dbTx *sqlx.Tx, id string) (err error) {
	defer trace.Default(&ctx, &err)()

	var querier db.Querier = dbTx
	if dbTx == nil {
		querier = c.db.GetMaster()
	}

	q := `UPDATE %s SET deleted_at=now() WHERE id = ?`
	q = fmt.Sprintf(q, c.tableName)

	_, err = querier.ExecContext(ctx, c.db.Rebind(q), id)
	if err != nil {
		return err
	}

	return nil
}

func (c *StdCRUD[T]) Get(ctx context.Context, dbTx *sqlx.Tx, param db.Params) (data []T, total int, err error) {
	defer trace.Default(&ctx, &err)()

	var querier db.Querier = dbTx
	if dbTx == nil {
		querier = c.db.GetSlave()
	}

	q := `SELECT id,%s FROM %s`
	q = fmt.Sprintf(q, c.getFieldsString(), c.tableName)

	countQ := `SELECT COUNT(1) FROM %s`
	countQ = fmt.Sprintf(countQ, c.tableName)

	param.Where = append(param.Where, db.Where{
		Field: "deleted_at",
		Op:    db.IsNullOp,
	})

	param.BuildWhere()
	countQ, countArgs := param.GetQuery(countQ)
	q, args := param.BuildSort().BuildPagination().GetQuery(q)

	err = querier.GetContext(ctx, &total, c.db.Rebind(countQ), countArgs...)
	if err != nil {
		return nil, 0, err
	}

	err = querier.SelectContext(ctx, &data, c.db.Rebind(q), args...)
	if err != nil {
		return data, 0, err
	}

	return data, total, nil
}

func (c *StdCRUD[T]) GetByIDs(ctx context.Context, dbTx *sqlx.Tx, ids []int, pg entity.Pagination) (data []T, total int, err error) {
	defer trace.Default(&ctx, &err)()

	data, total, err = c.Get(ctx, dbTx, db.Params{
		Where: []db.Where{
			{
				Field: "id",
				Op:    db.InOp,
				Value: ids,
			},
		},
		Pagination: pg.QBPaginate(),
		Sort:       pg.QBSort(),
	})
	if err != nil {
		return data, 0, err
	}

	return data, total, nil
}
