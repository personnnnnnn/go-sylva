package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sylva/sylva"
)

func main() {
	fmt.Println("Press enter to exit")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("$ ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			break
		}

		sylva.Execute(input, "<stdin>")
		fmt.Println()
	}
}
