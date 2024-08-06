package main

import (
	"bufio"
	"fmt"
	"os"
)

type command struct {
	syntax   string
	help     string
	function func() error
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("pokedex>")
		reader.Scan()
		println(reader.Text())
	}
}
