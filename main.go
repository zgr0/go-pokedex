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
func commandList() func() error {
	return func() error {
		for i:=1;i<=20;i++{
		url := fmt.Sprintf("https://pokeapi.co/api/v2/location/%d/?limit=6&offset=0",i)
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
	return nil 
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
	cmdMap["list"] = cliCommand{
		name:     "list",
		callback: commandList(),
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
