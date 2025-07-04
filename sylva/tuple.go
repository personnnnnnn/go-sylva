package sylva

import "fmt"

type Tuple struct {
	Keys   []string
	Values []Value
}

func (t *Tuple) GetIndex(index int) (Value, error) {
	if index < 0 || index >= len(t.Values) {
		return nil, fmt.Errorf(
			"index '%v' out of range for tuple of length '%v",
			index,
			len(t.Values),
		)
	}
	return t.Values[index], nil
}

func (t *Tuple) GetKey(key string) (Value, error) {
	for i, k := range t.Keys {
		if k == key {
			return t.Values[i], nil
		}
	}
	return nil, fmt.Errorf("tuple has no key '%v'", key)
}
