package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}
type Area struct {
	Name string `json:"name"`
}

type Data struct {
	Areas []Area `json:"areas"`
}
type config struct {
	start int
	end   int
}

var limit = config{start: 1, end: 5}

func handleInput(cmd cliCommand, ok bool) {
	if ok {
		err := cmd.callback()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	} else {
		fmt.Println("Unknown command")
	}
}

func commandHelp(cmdMap map[string]cliCommand) func() error {
	return func() error {
		fmt.Println("HELP!!!!")
		for _, v := range cmdMap {
			fmt.Printf("%s: %s\n", v.name, v.description)
		}
		return nil
	}
}

/*
	func getCommands() map[string]cliCommand {
		return map[string]cliCommand{
			"help": {
				name:        "help",
				description: "Displays a help message",
				callback:    commandHelp,
			},
			"exit": {
				name:        "exit",
				description: "Exit the Pokedex",
				callback:    commandExit,
			},
		}
	}
*/
func commandmap() func() error {
	return func() error {
		if limit.start <= 0 || limit.end < 6 {
			limit.start = 1
			limit.end = 5
		}
		for i := limit.start; i < limit.end; i++ {
			if i== 21 {
				continue
			}
			url := fmt.Sprintf("https://pokeapi.co/api/v2/location/%d/?limit=1&offset=0", i)
			res, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			body, err := io.ReadAll(res.Body)
			defer res.Body.Close()
			// Unmarshal JSON data
			var data Data
			jsonerr := json.Unmarshal(body, &data)
			if jsonerr != nil {
				fmt.Println("Error unmarshalling JSON:", jsonerr)
				limit.start = 1
				limit.end = 5
				fmt.Printf("i:%d", i)
				return nil
			}
			// Print the name of each area
			for _, area := range data.Areas {
				fmt.Println(area.Name)
			}

			if res.StatusCode > 299 {
				log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
			}
			if err != nil {
				log.Fatal(err)
			}
		}
		limit.start += 5
		limit.end += 5
		return nil
	}

}
func commandMapBack() func() error {
	return func() error {
		fmt.Printf("%d\n", limit.start)
		fmt.Printf("%d\n", limit.end)

		if limit.start <= 0 || limit.end < 6 {
			fmt.Println("You are at the starting point. You can't go back.")
			limit.start = 1
			limit.end = 5
			return nil
		} else {
			limit.start -= 20
			limit.end -= 20
			commandFunc := commandmap() // Get the function
			return commandFunc()        // Call the function
		}
	}
}

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
