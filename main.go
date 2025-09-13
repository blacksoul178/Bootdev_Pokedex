package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	lower_text := strings.ToLower(text)
	split_text := strings.Fields(lower_text)
	return split_text
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	// "help":{
	// 	name:  "help",
	// 	description: "returns help for available commands",
	// 	callback: commandHelp,

	// },
}

func main() {
	fmt.Print("welcome to the Boot.Dev Pokedex!, please enter a search.\n")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex>")
		scanner.Scan()
		text := scanner.Text()
		cleaned := cleanInput(text)
		switch cleaned[0] {
		case "exit":
			commands["exit"].callback()
		default:
			fmt.Print("Unknown Command\n")
		}

	}

}
