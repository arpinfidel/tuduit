package entity

import (
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

type Base36[T constraints.Unsigned] struct {
	V T
}

func NewBase36[T constraints.Unsigned](v T) Base36[T] {
	return Base36[T]{
		V: v,
	}
}

func NewBase36FromString[T constraints.Unsigned](s string) Base36[T] {
	return Base36[T]{
		V: T(0),
	}
}

type Base36Parser interface {
	ParseBase36(s string) error
}

func (b *Base36[T]) ParseBase36(s string) error {
	v, err := strconv.ParseUint(s, 36, 64)
	if err != nil {
		return err
	}

	b.V = T(v)
	return nil
}

func (b Base36[T]) String() string {
	return strings.ToUpper(strconv.FormatUint(uint64(b.V), 36))
}

func (b Base36[T]) MarshalYAML() (any, error) {
	return b.String(), nil
}

func (b *Base36[T]) UnmarshalYAML(unmarshal func(any) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	return b.ParseBase36(s)
}
