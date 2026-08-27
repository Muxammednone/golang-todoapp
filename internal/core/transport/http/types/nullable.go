package core_http_types

import (
	"encoding/json"

	"github.com/Muxammednone/golang-todoapp/internal/core/domain"
)

type Nullable[T any] struct {
	domain.Nullable[T]
}

func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	n.Set = true

	if string(b) == "null" {
		n.Value = nil
		return nil
	}

	var Value T
	if err := json.Unmarshal(b, &Value); err != nil {
		return err
	}
	n.Value = &Value
	return nil
}

func (n *Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{
		Set:   n.Set,
		Value: n.Value,
	}
}
