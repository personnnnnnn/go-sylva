package main

import (
	"fmt"
	"sylva/sylva"
)

func main() {
	m := make(map[sylva.Value]sylva.Value)
	m["x"] = true
	m["y"] = sylva.Tuple{
		Keys:   []string{"x", "y"},
		Values: []sylva.Value{"Hello, World!", false},
	}
	fmt.Println("Hello, World!", m["x"], m["y"])
}
