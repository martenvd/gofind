package app

import (
	"bufio"
	"fmt"
	"os"
)

func CacheDirs(homeDir string, dirs []string) {
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

func ReadDirs(homeDir string) ([]string, error) {
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
