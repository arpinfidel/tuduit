package entity

type Example struct {
	StdFields

	Name string `json:"name" db:"name"`
}
