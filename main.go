// GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o gofind main.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"golang.org/x/term"
)

func main() {

	updateCache := flag.Bool("u", false, "Whether or not to update the gofind cache.")
	updateCacheFullName := flag.Bool("update", false, "Whether or not to update the gofind cache.")
	flag.Parse()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	if *updateCache || *updateCacheFullName || !fileExists(homeDir+"/.gofind/dirs.txt") {
		dirs, err := walkPaths()
		if err != nil {
			log.Fatal(err)
		}
		cacheDirs(homeDir, dirs)
	}

	dirs, err := readDirs(homeDir)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(dirs)
	// prompt(dirs)
	fuzzyFind(dirs)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func isFlag() bool {
	flags := []string{"-u", "-update"}

	for _, flag := range flags {
		if os.Args[1] == flag {
			return true
		}
	}
	return false
}

func walkPaths() ([]string, error) {
	var dirs []string

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rgx, err := regexp.Compile(`\B\.git(\b|$|\s)`)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if info.IsDir() {
			files, err := os.ReadDir(path)
			if err != nil {
				fmt.Println(err)
				return err
			}

			for _, file := range files {
				if rgx.MatchString(file.Name()) {
					dirs = append(dirs, path)
					break
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return dirs, nil
}

func prompt(dirs []string) {

	var relevantDirs []string

	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Check if a command-line argument was provided
	if len(os.Args) < 2 {
		fmt.Println("Please provide a root directory for the search.")
		return
	}

	arg := 1
	if isFlag() {
		if len(os.Args) > 2 {
			arg = 2
		} else {
			fmt.Println("You have updated your cache!")
			return
		}
	}

	for _, dir := range dirs {
		if strings.Contains(dir, currentDir) && strings.Contains(strings.ToLower(dir), strings.ToLower(os.Args[arg])) {
			relevantDirs = append(relevantDirs, dir)
		}
	}

	prompt := promptui.Select{
		Label: "Select Directory",
		Items: relevantDirs,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("You selected %s\n", result)

	cmd := exec.Command("code", result)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

}

func cacheDirs(homeDir string, dirs []string) {
	// Create a file
	os.MkdirAll(homeDir+"/.gofind", os.ModePerm)
	file, err := os.Create(homeDir + "/.gofind/dirs.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write strings to the file
	for _, dir := range dirs {
		_, err := file.WriteString(dir + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
}

func readDirs(homeDir string) ([]string, error) {
	// Open the file for reading
	file, err := os.Open(homeDir + "/.gofind/dirs.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	// Scanner to read from file
	scanner := bufio.NewScanner(file)

	// Set to store strings
	var dirs []string

	// Read strings from the file
	for scanner.Scan() {
		dir := scanner.Text()
		dirs = append(dirs, dir)
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	return dirs, nil
}

func fuzzyFind(dirs []string) {
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
