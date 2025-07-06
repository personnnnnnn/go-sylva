package util

import "strings"

type Position struct {
	Col int
	Ln  int
	Idx int
}

// based of off https://github.com/davidcallanan/py-myopl-code/blob/master/ep14/strings_with_arrows.py
func StringsWithArrows(text string, fileName string, start, end Position, linePrefix string) string {
	result := ""

	idxStart := intMax(0, rfind(text, '\n', 0, start.Idx))
	idxEnd := find(text, '\n', idxStart+1, len(text))
	if idxEnd < 0 {
		idxEnd = len(text)
	}

	lineCount := end.Ln - start.Ln + 1
	for i := range lineCount {
		println(len(text), idxStart, idxEnd)
		line := text[idxStart:idxEnd]
		var colStart int
		if i == 0 {
			colStart = start.Col
		} else {
			colStart = 0
		}
		var colEnd int
		if i == lineCount-1 {
			colEnd = end.Col
		} else {
			colEnd = len(line)
		}

		result += line + "\n"
		result += strings.Repeat(" ", colStart) + strings.Repeat("^", colEnd-colStart)

		idxStart = idxEnd
		idxEnd = find(text, '\n', idxStart+1, len(text))
		if idxEnd < 0 {
			idxEnd = len(text)
		}
	}

	result = strings.ReplaceAll(result, "\t", "")
	result = linePrefix + result
	result = strings.ReplaceAll(result, "\n", "\n"+linePrefix)

	return result
}

func intMax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func rfind(text string, char byte, start, end int) int {
	for i := end - 1; i >= start; i-- {
		if text[i] == char {
			return i
		}
	}
	return -1
}

func find(text string, char byte, start, end int) int {
	for i := start; i < end; i++ {
		if text[i] == char {
			return i
		}
	}
	return -1
}
