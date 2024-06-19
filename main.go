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

	"github.com/lithammer/fuzzysearch/fuzzy"
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
	// Set terminal to raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	input := ""
	selectedIndex := 0
	fmt.Print("Search: ")

	for {
		char := make([]byte, 3)
		_, err := os.Stdin.Read(char)
		if err != nil {
			fmt.Println("Error reading character:", err)
			break
		}
		matches := fuzzy.Find(input, dirs)

		if char[0] == 27 && char[1] == 91 { // Arrow keys are escape sequences starting with 27, 91
			switch char[2] {
			case 65: // Up arrow
				if selectedIndex > 0 {
					selectedIndex--
				}
			case 66: // Down arrow
				selectedIndex++
			}
		} else {
			switch char[0] {
			case 3: // Ctrl+C
				return
			case 127: // Backspace
				if len(input) > 0 {
					input = input[:len(input)-1]
					selectedIndex = 0 // Reset selection on input change
				}
			case 13: // Enter
				if len(matches) > 0 && selectedIndex < len(matches) {
					fmt.Print("\033[H\033[2J")
					fmt.Print("\033[3K\rSelected directory:" + matches[selectedIndex] + "\r\n")
					return
				}
			default:
				input += string(char[0])
				selectedIndex = 0 // Reset selection on input change
			}
		}

		// Clear current line and reprint search prompt and input
		// fmt.Print("Search: ", input)

		// Perform fuzzy search and display results
		if selectedIndex >= len(matches) { // Ensure selectedIndex is within bounds
			selectedIndex = len(matches) - 1
		}

		// Clear screen
		// fmt.Print("\033[H\033[2J")

		for i, match := range matches {
			// fmt.Println("\033[3K\r", match)
			if i == selectedIndex {
				fmt.Println("\033[3K\r", "\033[1;4m"+match+"\033[0m") // Highlight selected match
			} else {
				fmt.Println("\033[3K\r", match)
			}
		}

		fmt.Print("\033[2K\rSearch: ", input) // Clear line and reprint input prompt and current input

	}
}
