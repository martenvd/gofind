package core

import (
	"fmt"
	"os"

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

	relevantDirs = utils.GetFilteredResults(currentDir, os.Args[arg], dirs)

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

	utils.OpenInVSCodeFromFinder(result, len(relevantDirs))
}
