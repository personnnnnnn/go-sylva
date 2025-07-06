package sylva

import "fmt"

func Add(a, b Value) (Value, error) {
	switch va := a.(type) {
	case int64:
		switch vb := b.(type) {
		case int64:
			return va + vb, nil
		case float64:
			return float64(va) + vb, nil
		default:
			return nil, fmt.Errorf("cannot add '%T' and '%T'", a, b)
		}
	case float64:
		switch vb := b.(type) {
		case int64:
			return va + float64(vb), nil
		case float64:
			return va + vb, nil
		default:
			return nil, fmt.Errorf("cannot add '%T' and '%T'", a, b)
		}
	default:
		return nil, fmt.Errorf("cannot add '%T' and '%T'", a, b)
	}
}

func Sub(a, b Value) (Value, error) {
	switch va := a.(type) {
	case int64:
		switch vb := b.(type) {
		case int64:
			return va - vb, nil
		case float64:
			return float64(va) - vb, nil
		default:
			return nil, fmt.Errorf("cannot subtract '%T' and '%T'", a, b)
		}
	case float64:
		switch vb := b.(type) {
		case int64:
			return va - float64(vb), nil
		case float64:
			return va - vb, nil
		default:
			return nil, fmt.Errorf("cannot subtract '%T' and '%T'", a, b)
		}
	default:
		return nil, fmt.Errorf("cannot subtract '%T' and '%T'", a, b)
	}
}

func Mul(a, b Value) (Value, error) {
	switch va := a.(type) {
	case int64:
		switch vb := b.(type) {
		case int64:
			return va * vb, nil
		case float64:
			return float64(va) * vb, nil
		default:
			return nil, fmt.Errorf("cannot multiply '%T' and '%T'", a, b)
		}
	case float64:
		switch vb := b.(type) {
		case int64:
			return va * float64(vb), nil
		case float64:
			return va * vb, nil
		default:
			return nil, fmt.Errorf("cannot multiply '%T' and '%T'", a, b)
		}
	default:
		return nil, fmt.Errorf("cannot multiply '%T' and '%T'", a, b)
	}
}

func Div(a, b Value) (Value, error) {
	switch va := a.(type) {
	case int64:
		switch vb := b.(type) {
		case int64:
			if vb == 0 {
				return nil, fmt.Errorf("cannot divide by 0")
			}
			return va / vb, nil
		case float64:
			if vb == 0 {
				return nil, fmt.Errorf("cannot divide by 0")
			}
			return float64(va) / vb, nil
		default:
			return nil, fmt.Errorf("cannot divide '%T' and '%T'", a, b)
		}
	case float64:
		switch vb := b.(type) {
		case int64:
			if vb == 0 {
				return nil, fmt.Errorf("cannot divide by 0")
			}
			return va / float64(vb), nil
		case float64:
			if vb == 0 {
				return nil, fmt.Errorf("cannot divide by 0")
			}
			return va / vb, nil
		default:
			return nil, fmt.Errorf("cannot divide '%T' and '%T'", a, b)
		}
	default:
		return nil, fmt.Errorf("cannot divide '%T' and '%T'", a, b)
	}
}

func Mod(a, b Value) (Value, error) {
	switch va := a.(type) {
	case int64:
		switch vb := b.(type) {
		case int64:
			if vb == 0 {
				return nil, fmt.Errorf("cannot mod by 0")
			}
			return va % vb, nil
		default:
			return nil, fmt.Errorf("cannot mod '%T' and '%T'", a, b)
		}
	default:
		return nil, fmt.Errorf("cannot mod '%T' and '%T'", a, b)
	}
}

func Umn(v Value) (Value, error) {
	switch x := v.(type) {
	case int64:
		return -x, nil
	case float64:
		return -x, nil
	default:
		return nil, fmt.Errorf("cannot apply umn to type '%T'", x)
	}
}

func Concat(a, b Value) string {
	return fmt.Sprintf("%v%v", a, b)
}
