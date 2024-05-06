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
)

func main() {

	updateCache := flag.Bool("u", false, "Whether or not to update the gofind cache.")
	flag.Parse()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	if *updateCache || !fileExists(homeDir+"/.gofind/dirs.txt") {
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

	prompt(dirs)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func isFlag() bool {
	flags := []string{"-u"}

	for _, flag := range flags {
		if os.Args[1] == flag {
			return true
		}
	}
	return false
}

func walkPaths() ([]string, error) {
	var dirs []string

	err := filepath.Walk("/home/martenvd", func(path string, info os.FileInfo, err error) error {
		if info != nil && info.IsDir() {
			output, err := exec.Command("ls", "-a", path).Output()
			if err != nil {
				fmt.Println(err)
				return err
			}

			git, err := regexp.MatchString(`\B\.git(\b|$|\s)`, string(output))
			if err != nil {
				fmt.Println(err)
			}
			if git {
				dirs = append(dirs, path)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return dirs, err
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
