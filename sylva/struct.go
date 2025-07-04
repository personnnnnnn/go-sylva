package sylva

import "fmt"

type Struct struct {
	Keys   []string
	Values []Value
}

func (s *Struct) GetKey(key string) (Value, error) {
	for i, k := range s.Keys {
		if k == key {
			return s.Values[i], nil
		}
	}
	return nil, fmt.Errorf("struct has no key '%v'", key)
}

func (s *Struct) SetKey(key string, value Value) error {
	for i, k := range s.Keys {
		if k == key {
			s.Values[i] = value
			return nil
		}
	}
	return fmt.Errorf("struct has no key '%v'", key)
}
