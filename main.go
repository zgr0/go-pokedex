package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)



func main() {
	reader := bufio.NewScanner(os.Stdin)
	cmdMap := make(map[string]cliCommand)

	cmdMap["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp(cmdMap),
	}
	cmdMap["exit"] = cliCommand{
		name:        "exit",
		description: "exits from app",
	}
	cmdMap["map"] = cliCommand{
		name:     "map",
		callback: commandmap(),
	}
	cmdMap["mapb"] = cliCommand{
		name:     "map",
		callback: commandMapBack(),
	}
	for {
		fmt.Print("pokedex> ")
		if reader.Scan() {
			input := reader.Text()
			if strings.ToLower(input) == "exit" {
				break
			}

			cmd, ok := cmdMap[input]
			handleInput(cmd, ok)
		} else {
			fmt.Println("Error reading input")
			break
		}
	}
}
