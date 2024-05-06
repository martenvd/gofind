// GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o gofind main.go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
)

func main() {
	var dirs []string

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
	err = filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if info != nil && info.IsDir() && strings.Contains(strings.ToLower(path), strings.ToLower(os.Args[1])) {
			output, err := exec.Command("ls", "-a", path).Output()
			if err != nil {
				fmt.Println(err)
				return err
			}

			git, err := regexp.MatchString("\\.git\\s", string(output))
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
		return
	}

	prompt := promptui.Select{
		Label: "Select Directory",
		Items: dirs,
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
