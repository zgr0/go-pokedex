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
	cache      *cacheList
}

// NewClient -
func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: &cacheList{cacheMap: make(map[string]cache)},
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func main() {
	pokeClient := NewClient(5 * time.Second)
	cfg := &config{
		pokeapiClient: pokeClient,
	}
	reader := bufio.NewScanner(os.Stdin)
	fmt.Println("                                  ,'\\")
	fmt.Println("    _.----.        ____         ,'  _\\   ___    ___     ____")
	fmt.Println("_,-'       `.     |    |  /`.   \\,-'    |   \\  /   |   |    \\  |`.")
	fmt.Println("\\      __    \\    '-.  | /   `.  ___    |    \\/    |   '-.   \\ |  |")
	fmt.Println(" \\.    \\ \\   |  __  |  |/    ,','_  `.  |          | __  |    \\|  |")
	fmt.Println("   \\    \\/   /,' _`.|      ,' / / / /   |          ,' _`.|     |  |")
	fmt.Println("    \\     ,-'/  /   \\    ,'   | \\/ / ,`.|         /  /   \\  |     |")
	fmt.Println("     \\    \\ |   \\_/  |   `-.  \\    `'  /|  |    ||   \\_/  | |\\    |")
	fmt.Println("      \\    \\ \\      /       `-.`.___,-' |  |\\  /| \\      /  | |   |")
	fmt.Println("       \\    \\ `.__,'|  |`-._    `|      |__| \\/ |  `.__,'|  | |   |")
	fmt.Println("        \\_.-'       |__|    `-._ |              '-.|     '-.| |   |")
	for {
		fmt.Print("pokedex> ")
		reader.Scan()
		input := cleanInput(reader.Text())
		if len(input) == 0 {
			continue
		}
		arg := []string{}
		if len(input) > 1 {
			arg = input[1:]
		}
		commandname := input[0]
		fmt.Println(input[0])
		fmt.Println(arg)
		command, ok := getCommands()[commandname]
		if ok {
			err := command.callback(cfg, arg...)
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
