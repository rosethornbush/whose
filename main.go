package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: whose <query>")
		os.Exit(2)
	}

	query := os.Args[1]

	fmt.Println(query)
}
