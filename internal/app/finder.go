package app

import (
	"bufio"
	"fmt"
	"os"

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
	fmt.Println("Getypte letters: ")
	fmt.Print("\rZoek: ")

	reader := bufio.NewReader(os.Stdin)
	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			break
		}

		if char == '\033' { // Escape karakter gedetecteerd
			// Lees de volgende twee karakters om te bepalen of het "up" of "down" is
			nextChar, _ := reader.Peek(2)
			if string(nextChar) == "[A" { // "up" toets
				// Consumeer de karakters zodat ze niet in de stdout verschijnen
				reader.Discard(2)
				// Voeg hier eventueel logica toe om iets te doen wanneer "up" wordt gedrukt
				continue
			} else if string(nextChar) == "[B" { // "down" toets
				// Consumeer de karakters zodat ze niet in de stdout verschijnen
				reader.Discard(2)
				// Voeg hier eventueel logica toe om iets te doen wanneer "down" wordt gedrukt
				continue
			}
			// Voeg extra cases toe voor andere toetsen zoals "right" ([C) en "left" ([D) indien nodig
		}

		switch char {
		case '\r': // Enter key
			fmt.Print("\n\rZoekopdracht voltooid.\n\r")
			return
		case 127: // Backspace key
			if len(input) > 0 {
				input = input[:len(input)-1]
				// Wis de huidige regel en toon de bijgewerkte input, voeg een extra spatie toe om overgebleven karakters te overschrijven
				fmt.Print("\033[1A\033[2K\rGetypte letters: ", input, " \n\033[K\rZoek: ", input, " ")
				// Beweeg de cursor één positie naar links om de extra spatie niet als deel van de input te tonen
				fmt.Print("\033[1D")
			}
		default:
			input += string(char)
			// Update de getypte letters en het zoekveld
			fmt.Print("\033[1A\033[2K\rGetypte letters: ", input, "\n\033[K\rZoek: ", input)
		}
	}
}
