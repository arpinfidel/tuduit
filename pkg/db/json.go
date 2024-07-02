package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

var (
	_ sql.Scanner   = &JSON[any]{}
	_ driver.Valuer = JSON[any]{}

	_ json.Marshaler   = JSON[any]{}
	_ json.Unmarshaler = &JSON[any]{}
)

// JSON is a struct that scans and values to json
type JSON[T any] struct {
	V T
}

// Scan implements the sql.Scanner interface.
func (j *JSON[T]) Scan(value any) error {
	if value == nil {
		var zero T
		j.V = zero
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid type %T", value)
	}

	return json.Unmarshal(b, &j.V)
}

// Value implements the driver Valuer interface.
func (j JSON[T]) Value() (driver.Value, error) {
	b, err := json.Marshal(j.V)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalJSON implements the json.Marshaler interface.
func (j JSON[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.V)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (j *JSON[T]) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &j.V)
}

var (
	_ sql.Scanner   = &NullJSON[any]{}
	_ driver.Valuer = NullJSON[any]{}

	_ json.Marshaler   = NullJSON[any]{}
	_ json.Unmarshaler = &NullJSON[any]{}
)

// NullJSON is a struct that scans and values to json
type NullJSON[T any] struct {
	Valid bool
	V     T
}

// Scan implements the sql.Scanner interface.
func (j *NullJSON[T]) Scan(value any) error {
	if value == nil {
		var zero T
		j.Valid = false
		j.V = zero
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid type %T", value)
	}

	err := json.Unmarshal(b, &j.V)
	if err != nil {
		return err
	}

	j.Valid = true

	return nil
}

// Value implements the driver Valuer interface.
func (j NullJSON[T]) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}

	b, err := json.Marshal(j.V)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalJSON implements the json.Marshaler interface.
func (j NullJSON[T]) MarshalJSON() ([]byte, error) {
	if !j.Valid {
		return []byte(`null`), nil
	}

	return json.Marshal(j.V)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (j *NullJSON[T]) UnmarshalJSON(b []byte) error {
	if j == nil {
		return fmt.Errorf("jsondb: UnmarshalJSON on nil pointer")
	}

	if len(b) == 0 || string(b) == `null` {
		var zero T
		j.Valid = false
		j.V = zero
		return nil
	}

	err := json.Unmarshal(b, &j.V)
	if err != nil {
		return err
	}

	j.Valid = true
	return nil
}
