package entity

import "time"

type StdFields struct {
	ID        int64      `json:"id"         db:"id"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	CreatedBy string     `json:"created_by" db:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy string     `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

func (std StdFields) GetStdFields() StdFields {
	return std
}
