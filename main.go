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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	const initialURL = "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	url := initialURL
	var pokedexResponseModel PokedexResponseModel

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

			for _, result := range pokedexResponseModel.Results {
				fmt.Println(result.Name)
			}

		case "mapb":
			if url == initialURL {
				fmt.Println("you are already on the first page")
				break
			}
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

			for _, result := range pokedexResponseModel.Results {
				fmt.Println(result.Name)
			}
		}
	}
}
