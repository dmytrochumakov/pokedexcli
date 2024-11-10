package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dmytrochumakov/pokedexcli/internal/pokeapi"
)

type config struct {
	apiClient        pokeapi.Client
	nextLocationsURL *string
	prevLocationsURL *string
	caughtPokemon    map[string]pokeapi.Pokemon
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	client := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	cfg := &config{
		caughtPokemon: map[string]pokeapi.Pokemon{},
		apiClient:     client,
	}

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		args := strings.Fields(input)

		if len(args) < 1 {
			continue
		}

		switch args[0] {
		case "help":
			fmt.Println("Welcome to the Pokedex!")
			fmt.Println("Usage:")
			fmt.Println("help: Displays a help message")
			fmt.Println("exit: Exit the Pokedex")

		case "exit":
			return

		case "map":
			locationsResp, err := cfg.apiClient.ListLocations(cfg.nextLocationsURL)
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, locatiun := range locationsResp.Results {
				fmt.Println(locatiun.Name)
			}

		case "mapb":
			if cfg.prevLocationsURL == nil {
				fmt.Println("you are alredy on first page")
			}

			locationsResp, err := cfg.apiClient.ListLocations(cfg.nextLocationsURL)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, location := range locationsResp.Results {
				fmt.Println(location.Name)
			}

		case "explore":
			if len(args) < 2 {
				fmt.Println("Usage: explore <text>")
			} else {
				pokemonName := args[1]
				fmt.Println("pokemon nam: " + pokemonName)
				location, err := cfg.apiClient.GetLocation(pokemonName)
				if err != nil {
					fmt.Println(err)
					return
				}
				for _, pokemon := range location.PokemonEncounters {
					fmt.Println(" - " + pokemon.Pokemon.Name)
				}
			}

		case "catch":
			if len(args) < 2 {
				fmt.Println("Usage: catch <text>")
			} else {
				pokemonName := args[1]
				pokemon, err := cfg.apiClient.GetPokemon(pokemonName)
				if err != nil {
					fmt.Println(err)
					return
				}

				res := rand.Intn(pokemon.BaseExperience)

				fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

				if res > 40 {
					fmt.Printf("%s escaped!\n", pokemon.Name)
					break
				}

				fmt.Printf("%s was caught!\n", pokemon.Name)
				cfg.caughtPokemon[pokemon.Name] = pokemon
			}

		case "inspect":
			if len(args) < 2 {
				fmt.Println("Usage: inspect <text>")
			} else {
				pokemonName := args[1]
				pokemon, ok := cfg.caughtPokemon[pokemonName]
				if !ok {
					fmt.Println("you have no caught pokemons with that name")
					break
				}

				fmt.Println("Name:", pokemon.Name)
				fmt.Println("Height:", pokemon.Height)
				fmt.Println("Weight:", pokemon.Weight)
				fmt.Println("Stats:")
				for _, stat := range pokemon.Stats {
					fmt.Printf("  -%s: %v\n", stat.Stat.Name, stat.BaseStat)
				}
				fmt.Println("Types:")
				for _, typeInfo := range pokemon.Types {
					fmt.Println("  -", typeInfo.Type.Name)
				}
			}
		}
	}
}
