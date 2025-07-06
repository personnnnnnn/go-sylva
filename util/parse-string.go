package util

import (
	"fmt"
	"strconv"
)

func ParseString(s string) (string, error) {
	res := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' {
			res = append(res, s[i])
			continue
		}
		i++
		if i >= len(s) {
			return "", fmt.Errorf("incomplete escape sequence at end of string")
		}
		switch s[i] {
		case '\\':
			res = append(res, '\\')
		case '\'':
			res = append(res, '\'')
		case '"':
			res = append(res, '"')
		case 'n':
			res = append(res, '\n')
		case 'r':
			res = append(res, '\r')
		case 't':
			res = append(res, '\t')
		case 'b':
			res = append(res, '\b')
		case 'f':
			res = append(res, '\f')
		case 'a':
			res = append(res, '\a')
		case 'v':
			res = append(res, '\v')
		case 'x': // hex
			if i+2 >= len(s) {
				return "", fmt.Errorf("incomplete hex escape sequence at index %d", i)
			}
			hexDigits := s[i+1 : i+3]
			val, err := strconv.ParseUint(hexDigits, 16, 8)
			if err != nil {
				return "", fmt.Errorf("invalid hex escape \\x%s at index %d", hexDigits, i)
			}
			res = append(res, byte(val))
			i += 2
		default:
			if s[i] >= '0' && s[i] <= '7' { // octal
				start := i
				end := i + 1
				for end < len(s) && end-start < 3 && s[end] >= '0' && s[end] <= '7' {
					end++
				}
				octalDigits := s[start:end]
				val, err := strconv.ParseUint(octalDigits, 8, 8)
				if err != nil {
					return "", fmt.Errorf("invalid octal escape \\%s", octalDigits)
				}
				res = append(res, byte(val))
				i = end - 1
			} else {
				return "", fmt.Errorf("unknown escape sequence: \\%c", s[i])
			}
		}
	}
	return string(res), nil
}
