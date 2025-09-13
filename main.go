package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type config struct {
	Next_url     string
	Previous_url string
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands map[string]cliCommand

type LocationArea struct {
	Name string `json:"name"`
}

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

func commandHelp() error {
	fmt.Println("Usage:")
	for name, cmd := range commands {
		fmt.Printf("%s : %s\n", name, cmd.description)
	}
	return nil
}

func commandMap() error {
	for i := 1; i <= 20; i++ {
		url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d", i)
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", url, err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}
		area := LocationArea{}
		if err := json.Unmarshal(body, &area); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
		fmt.Printf("%s\n", area.Name)

	}
	return nil
}

func main() {
	fmt.Println("Welcome to the Pokedex!")
	scanner := bufio.NewScanner(os.Stdin)

	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "help",
			description: "display names of map locations",
			callback:    commandMap,
		},
	}

	for {
		fmt.Print("Pokedex>")
		scanner.Scan()
		text := scanner.Text()
		cleaned := cleanInput(text)
		switch cleaned[0] {
		case "exit":
			commands["exit"].callback()
		case "help":
			commands["help"].callback()
		case "map":
			commands["map"].callback()
		default:
			fmt.Print("Unknown Command\n")
		}

	}

}
