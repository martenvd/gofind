package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func FuzzyFind(dirs []string) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Failed to set terminal to raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	input := ""
	fmt.Print("\rSearch: ")

	// selectedIndex := 0

	reader := bufio.NewReader(os.Stdin)
	for {
		results := getFilteredResults(input, dirs)
		for i, result := range results {
			if len(results) != 0 {
				if i == 0 {
					fmt.Print("\r", result)
				} else {
					fmt.Print("\n\r", result)
				}
			}
		}
		if (len(results) != 0 && input != "") || (len(results) != 0 && input == "") {
			fmt.Print("\n\rSearch: ", input)
		} else {
			fmt.Print("\rSearch: ", input)
		}

		char, _, err := reader.ReadRune()
		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			break
		}

		switch char {
		case '\033': // Escape karakter
			nextChar, _ := reader.Peek(2)
			if string(nextChar) == "[A" { // "up" toets
				// Consumeer de karakters zodat ze niet in de stdout verschijnen
				reader.Discard(2)
				// Voeg hier eventueel logica toe om iets te doen wanneer "up" wordt gedrukt

				// fmt.Print("\rSearch: ", input)

				// if selectedIndex > 0 {
				// 	selectedIndex--
				// }
				continue
			} else if string(nextChar) == "[B" { // "down" toets
				// Consumeer de karakters zodat ze niet in de stdout verschijnen
				reader.Discard(2)
				// Voeg hier eventueel logica toe om iets te doen wanneer "down" wordt gedrukt

				// fmt.Print("\rSearch: ", input)
				// if selectedIndex < len(dirs)-1 {
				// }
				continue
			}
			// Voeg extra cases toe voor andere toetsen zoals "right" ([C) en "left" ([D) indien nodig
		case '\r': // Enter key
			fmt.Print("\n\rSearch complete\n\r")
			return
		case 127: // Backspace key
			if len(input) > 0 {
				// Clear screen
				fmt.Print("\033[H\033[2J")
				input = input[:len(input)-1]

				// Wis de huidige regel en toon de bijgewerkte input, voeg een extra spatie toe om overgebleven karakters te overschrijven
				// fmt.Print("\n\033[K\rSearch: ", input, " ")
				// if len(results) != 0 {
				// 	fmt.Print("\n\rSearch: ", input)
				// } else {
				// 	fmt.Print("\rSearch: ", input)
				// }
				// Beweeg de cursor één positie naar links om de extra spatie niet als deel van de input te tonen
				// fmt.Print("\033[1D")
			}
		default:
			// Clear screen
			fmt.Print("\033[H\033[2J")
			input += string(char)

			// fmt.Print("\n\033[K\rSearch: ", input)
			// fmt.Print("\n\rSearch: ", input)
		}
	}
}

func getFilteredResults(input string, dirs []string) []string {
	filteredResults := []string{}
	for _, dir := range dirs {
		if strings.Contains(strings.ToLower(dir), strings.ToLower(input)) {
			filteredResults = append(filteredResults, dir)
		}
	}
	return filteredResults
}
