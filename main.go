package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type config struct {
	pokeapiClient    Client
	nextLocationsURL *string
	prevLocationsURL *string
}

// Client -
type Client struct {
	httpClient http.Client
}

// NewClient -
func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func main() {
	pokeClient := NewClient(5 * time.Second)
	cfg := &config{
		pokeapiClient: pokeClient,
	}
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("pokedex> ")
		reader.Scan()
		input := cleanInput(reader.Text())
		if len(input) == 0 {
			continue
		}
		commandname := input[0]
		fmt.Println(input[0])
		command, ok := getCommands()[commandname]
		if ok {
			err := command.callback(cfg)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}

}
func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Get the next page of locations",
			callback:    commandMapf,
		},
		"mapb": {
			name:        "mapb",
			description: "Get the previous page of locations",
			callback:    commandMapb,
		},
	}
}
