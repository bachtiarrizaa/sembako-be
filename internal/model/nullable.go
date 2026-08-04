package model

import (
	"bytes"
	"encoding/json"
)

type Nullable[T any] struct {
	Value *T
	IsSet bool
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	n.IsSet = true

	if bytes.Equal(data, []byte("null")) {
		n.Value = nil
		return nil
	}

	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	n.Value = &v
	return nil
}
