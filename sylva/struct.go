package sylva

import "fmt"

type Struct struct {
	Keys   []string
	Values []Value
}

func (t *Struct) GetKey(key string) (Value, error) {
	for i, k := range t.Keys {
		if k == key {
			return t.Values[i], nil
		}
	}
	return nil, fmt.Errorf("tuple has no key '%v'", key)
}

func (t *Struct) SetKey(key string, value Value) error {
	for i, k := range t.Keys {
		if k == key {
			t.Values[i] = value
			return nil
		}
	}
	return fmt.Errorf("tuple has no key '%v'", key)
}
