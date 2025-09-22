package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
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
	callback    func(args []string) error
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
type exploreAreaResults struct {
	PokemonEncounters []PokemonEncounters `json:"pokemon_encounters"`
}
type PokemonEncounters struct {
	Pokemon Pokemon `json:"pokemon"`
}
type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type PokemonCatching struct {
	Name    string         `json:"name"`
	BaseExp int            `json:"base_experience"`
	Height  int            `json:"height"`
	Weight  int            `json:"weight"`
	Stats   []PokemonStats `json:"stats"`
	Types   []Types        `json:"types"`
}
type Stat struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type PokemonStats struct {
	BaseStat int  `json:"base_stat"`
	Effort   int  `json:"effort"`
	Stat     Stat `json:"stat"`
}
type Type struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Types struct {
	Slot int  `json:"slot"`
	Type Type `json:"type"`
}
type Pokedex map[string]PokemonCatching

func cleanInput(text string) []string {
	lower_text := strings.ToLower(text)
	split_text := strings.Fields(lower_text)
	return split_text
}
func commandExit([]string) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}
func commandHelp([]string) error {
	fmt.Println("Usage:")
	for name, cmd := range commands {
		fmt.Printf("%s : %s\n", name, cmd.description)
	}
	return nil
}
func commandMap(cfg *config, cache *pokecache.Cache) func([]string) error {
	return func([]string) error {
		if cfg.NextURL == "" {
			fmt.Println("you have reached the last page: use mapb to navigate to previous page")
			return nil
		}

		var body []byte
		cvalue, cbool := cache.Get(cfg.NextURL)
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
func commandMapBack(cfg *config, cache *pokecache.Cache) func([]string) error {
	return func([]string) error {
		if cfg.PreviousURL == "" {
			fmt.Println("you are on page one")
			return nil
		}

		var body []byte
		cvalue, cbool := cache.Get(cfg.PreviousURL)
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

		page := LocationAreaList{}
		if err := json.Unmarshal(body, &page); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
		cfg.NextURL = page.Next
		cfg.PreviousURL = page.Previous
		for i := range page.Results {
			fmt.Printf("%s\n", page.Results[i].Name)
		}

		return nil
	}
}
func commandExploreLocation(cache *pokecache.Cache) func([]string) error {
	return func(exploreLocation []string) error {
		if len(exploreLocation) < 1 {
			fmt.Println("Please supply an area to explore")
			return nil
		}
		if len(exploreLocation) > 1 {
			fmt.Println("Please supply only ONE area to explore at a time")
			return nil
		}

		fmt.Printf("Exploring %s...\n", exploreLocation[0])
		var body []byte
		locationURL := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", exploreLocation[0])

		cvalue, cbool := cache.Get(locationURL)
		if cbool && len(cvalue) > 0 {
			body = cvalue
		} else {
			resp, err := http.Get(locationURL)
			if err != nil {
				return fmt.Errorf("failed to fetch %s: %w", locationURL, err)
			}
			defer resp.Body.Close()

			body, err = io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response: %w", err)
			}
			cache.Add(locationURL, body)
		}

		exploreResults := exploreAreaResults{}
		if err := json.Unmarshal(body, &exploreResults); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
		fmt.Println("Found Pokemon:")
		for _, r := range exploreResults.PokemonEncounters {
			fmt.Printf("-%s\n", r.Pokemon.Name)
		}

		return nil
	}
}
func commandCatch(cache *pokecache.Cache, pokedex Pokedex) func([]string) error {
	return func(pokemonToCatch []string) error {
		if len(pokemonToCatch) < 1 {
			fmt.Println("Please supply which pokemon you wish to catch")
			return nil
		}
		if len(pokemonToCatch) > 1 {
			fmt.Println("Please supply only ONE pokemon to catch at a time")
			return nil
		}
		name := pokemonToCatch[0]

		fmt.Printf("Throwing a Pokeball at %s...\n", name)
		var body []byte
		catchURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)

		cvalue, cbool := cache.Get(catchURL)
		if cbool && len(cvalue) > 0 {
			body = cvalue
		} else {
			resp, err := http.Get(catchURL)
			if err != nil {
				return fmt.Errorf("failed to fetch %s: %w", catchURL, err)
			}
			defer resp.Body.Close()

			body, err = io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response: %w", err)
			}
			cache.Add(catchURL, body)
		}

		catchResults := PokemonCatching{}
		if err := json.Unmarshal(body, &catchResults); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}

		roll := rand.Intn(100)
		var scaler int = catchResults.BaseExp / 100
		var toughness int = catchResults.BaseExp / (10 + scaler)
		catchDifficulty := 40
		if roll > (catchDifficulty + toughness) {
			if _, ok := pokedex[name]; !ok {
				pokedex[name] = catchResults
				fmt.Printf("%s was caught!\n", pokemonToCatch[0])
			} else {
				fmt.Printf("%s already in pokedex\n", name)
			}

		} else {
			fmt.Printf("%s escaped!\n", pokemonToCatch[0])
		}
		fmt.Printf("Roll : %v, diff : %v\n", roll, (catchDifficulty + toughness))

		return nil
	}
}
func commandInspect(cache *pokecache.Cache, pokedex Pokedex) func([]string) error {
	return func(pokemonToInspect []string) error {
		if len(pokemonToInspect) < 1 {
			fmt.Println("Please supply which pokemon you wish to inspect")
			return nil
		}
		if len(pokemonToInspect) > 1 {
			fmt.Println("Please supply only ONE pokemon to inspect at a time")
			return nil
		}
		name := pokemonToInspect[0]
		if _, ok := pokedex[name]; !ok {
			fmt.Printf("you have not caught a %s yet", name)
		} else {
			fmt.Printf("Name: %s\n", pokedex[name].Name)
			fmt.Printf("Height %v\n", pokedex[name].Height)
			fmt.Printf("Weight: %v\n", pokedex[name].Weight)
			fmt.Println("Stats:")
			for _, s := range pokedex[name].Stats {
				fmt.Printf("-%v: %v\n", s.BaseStat, s.Stat.Name)
			}
			fmt.Println("Type(s):")
			for _, t := range pokedex[name].Types {
				fmt.Printf("- %v\n", t.Type.Name)
			}
		}
		return nil
	}

}
func main() {
	fmt.Println("Welcome to the Pokedex!")
	scanner := bufio.NewScanner(os.Stdin)
	cacheInterval := 30 * time.Second
	cache := pokecache.NewCache(cacheInterval)

	pokedex := make(Pokedex)

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
		"explore": {
			name:        "explore",
			description: "Explore a location to list all pokemon located there, usage: explore [area]",
			callback:    commandExploreLocation(cache),
		},
		"catch": {
			name:        "catch",
			description: "catch a pokemon, usage: catch [pokemon name]",
			callback:    commandCatch(cache, pokedex),
		},
		"inspect": {
			name:        "inspect",
			description: "inspect the stats of a caught pokemon, use inspect [pokemon name]",
			callback:    commandInspect(cache, pokedex),
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
		cmd := cleaned[0]
		args := cleaned[1:]
		if c, ok := commands[cmd]; ok {
			if err := c.callback(args); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown Command")

		}
	}

	// switch cleaned[0] {
	// case "exit":
	// 	commands["exit"].callback(cleaned[1:])
	// case "help":
	// 	commands["help"].callback(cleaned[1:])
	// case "map":
	// 	commands["map"].callback(cleaned[1:])
	// case "mapb":
	// 	commands["mapb"].callback(cleaned[1:])
	// case "explore":
	// 	commands["explore"].callback(cleaned[1:])
	// default:
	// 	fmt.Print("Unknown Command\n")
	// }

}
