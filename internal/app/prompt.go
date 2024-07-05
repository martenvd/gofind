package app

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/martenvd/gofind/internal/utils"
)

func Prompt(dirs []string) {

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
	if utils.IsFlag() {
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
		Label:        "Select Directory",
		Items:        relevantDirs,
		HideSelected: true,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Println(err)
		return
	}

	currentItemName := strings.Split(result, "/")[len(strings.Split(result, "/"))-1]
	fmt.Println("To open the directory type:")
	fmt.Println()
	fmt.Print("\tcd ", result, "\n")
	fmt.Println()
	fmt.Printf("Opening: %s\n", currentItemName)

	cmd := exec.Command("code", result)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

}
