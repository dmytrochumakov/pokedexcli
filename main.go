package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type PokedexResponseModel struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type PokemonReponseModel struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	const initialURL = "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	url := initialURL
	var pokedexResponseModel PokedexResponseModel
	cache := NewCache(5)

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
			if pokedexResponseModel.Next != "" {
				url = pokedexResponseModel.Next
			}
			if dataFromCaching, exists := cache.Get(url); exists {
				if err := json.Unmarshal(dataFromCaching, &pokedexResponseModel); err != nil {
					fmt.Println(err)
					return
				}

			} else {
				resp, err := http.Get(url)

				if err != nil {
					fmt.Println(err)
					return
				}
				defer resp.Body.Close()

				decoder := json.NewDecoder(resp.Body)
				err = decoder.Decode(&pokedexResponseModel)

				if err != nil {
					fmt.Println(err)
					return
				}
			}

			dataForCaching, err := json.Marshal(pokedexResponseModel)
			if err != nil {
				fmt.Println(err)
				return
			}

			cache.Add(url, dataForCaching)

			for _, result := range pokedexResponseModel.Results {
				fmt.Println(result.Name)
			}

		case "mapb":
			if url == initialURL {
				fmt.Println("you are already on the first page")
				break
			}
			if dataFromCaching, exists := cache.Get(url); exists {
				if err := json.Unmarshal(dataFromCaching, &pokedexResponseModel); err != nil {
					fmt.Println(err)
					return
				}

			} else {
				url = pokedexResponseModel.Previous
				resp, err := http.Get(url)

				if err != nil {
					fmt.Println(err)
					return
				}
				defer resp.Body.Close()

				decoder := json.NewDecoder(resp.Body)
				err = decoder.Decode(&pokedexResponseModel)

				if err != nil {
					fmt.Println(err)
					return
				}
			}

			for _, result := range pokedexResponseModel.Results {
				fmt.Println(result.Name)
			}

		case "explore":
			if len(args) < 2 {
				fmt.Println("Usage: explore <text>")
			} else {
				area := args[1]
				fmt.Println("Exploring" + area + "...")
				fmt.Println("Found Pokemon:")

				if dataFromCaching, exists := cache.Get(area); exists {
					if err := json.Unmarshal(dataFromCaching, &pokedexResponseModel); err != nil {
						fmt.Println(err)
						return
					}

				} else {
					resp, err := http.Get("https://pokeapi.co/api/v2/location-area/" + area + "?offset=0&limit=20")

					if err != nil {
						fmt.Println(err)
						return
					}
					defer resp.Body.Close()

					var pokemonReponseModel PokemonReponseModel
					decoder := json.NewDecoder(resp.Body)
					err = decoder.Decode(&pokemonReponseModel)

					if err != nil {
						fmt.Println(err)
						return
					}

					dataForCaching, err := json.Marshal(pokemonReponseModel)
					if err != nil {
						fmt.Println(err)
						return
					}

					cache.Add(area, dataForCaching)

					for _, pokemon := range pokemonReponseModel.PokemonEncounters {
						fmt.Println(" - " + pokemon.Pokemon.Name)
					}
				}
			}
		}
	}
}
