package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/blacksoul178/Bootdev_Pokedex/internal/pokecache"
)

type config struct {
	NextURL     string
	PreviousURL string
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands map[string]cliCommand

type LocationAreaList struct {
	Results  []LocationAreaResults `json:"results"`
	Next     string                `json:"next"`
	Previous string                `json:"previous"`
}

type LocationAreaResults struct {
	Name string `json:"name"`
	URL  string `json:"url"`
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

func commandMap(cfg *config, cache *pokecache.Cache) func() error {
	return func() error {
		if cfg.NextURL == "" {
			fmt.Println("you have reached the last page: use mapb to navigate to previous page")
			return nil
		}

		var body []byte
		cvalue, cbool := cache.Get(cfg.NextURL)
		fmt.Printf("Cache found: %v, Value length: %d\n", cbool, len(cvalue))
		if cbool && len(cvalue) > 0 {
			body = cvalue
		} else {
			resp, err := http.Get(cfg.NextURL)
			if err != nil {
				return fmt.Errorf("failed to fetch %s: %w", cfg.NextURL, err)
			}
			defer resp.Body.Close()

			body, err = io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response: %w", err)
			}
			cache.Add(cfg.NextURL, body)
		}

		page := LocationAreaList{}
		if err := json.Unmarshal(body, &page); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
		cfg.NextURL = page.Next
		cfg.PreviousURL = page.Previous
		for _, r := range page.Results {
			fmt.Printf("%s\n", r.Name)
		}

		return nil
	}
}

func commandMapBack(cfg *config, cache *pokecache.Cache) func() error {
	return func() error {
		if cfg.PreviousURL == "" {
			fmt.Println("you are on page one")
			return nil
		}

		var body []byte
		cvalue, cbool := cache.Get(cfg.PreviousURL)
		fmt.Printf("Cache found: %v, Value length: %d\n", cbool, len(cvalue))
		if cbool && len(cvalue) > 0 {
			body = cvalue
		} else {
			resp, err := http.Get(cfg.PreviousURL)
			if err != nil {
				return fmt.Errorf("failed to fetch %s: %w", cfg.PreviousURL, err)
			}
			defer resp.Body.Close()

			body, err = io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response: %w", err)
			}
			cache.Add(cfg.PreviousURL, body)
		}

		area := LocationAreaList{}
		if err := json.Unmarshal(body, &area); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
		cfg.NextURL = area.Next
		cfg.PreviousURL = area.Previous
		for i := range area.Results {
			fmt.Printf("%s\n", area.Results[i].Name)
		}

		return nil
	}
}

func main() {
	fmt.Println("Welcome to the Pokedex!")
	scanner := bufio.NewScanner(os.Stdin)
	cacheInterval := 10 * time.Second
	cache := pokecache.NewCache(cacheInterval)

	cfg := config{
		NextURL: "https://pokeapi.co/api/v2/location-area/",
	}
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
			name:        "map",
			description: "display the next 20 names of map locations",
			callback:    commandMap(&cfg, cache),
		},
		"mapb": {
			name:        "mapb",
			description: "display the previous 20 names of map locations",
			callback:    commandMapBack(&cfg, cache),
		},
	}

	for {
		fmt.Print("Pokedex>")
		scanner.Scan()
		text := scanner.Text()
		cleaned := cleanInput(text)
		if len(cleaned) == 0 {
			continue
		}
		switch cleaned[0] {
		case "exit":
			commands["exit"].callback()
		case "help":
			commands["help"].callback()
		case "map":
			commands["map"].callback()
		case "mapb":
			commands["mapb"].callback()
		default:
			fmt.Print("Unknown Command\n")
		}

	}

}
