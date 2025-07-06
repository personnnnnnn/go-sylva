package main

import "sylva/sylva"

func main() {
	data := "items = [1, 2] res = items[2]"
	sylva.Execute(data, "<stdin>")
}
