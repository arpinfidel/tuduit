package entity

import "time"

type StdFields struct {
	ID        int64      `db:"id"         json:"id"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	CreatedBy int64      `db:"created_by" json:"created_by"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	UpdatedBy int64      `db:"updated_by" json:"updated_by"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	DeletedBy *int64     `db:"deleted_by" json:"deleted_by"`
}

func (std StdFields) GetStdFields() StdFields {
	return std
}
