package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func IsFlag() bool {
	flags := []string{"-u", "-update"}

	for _, flag := range flags {
		if os.Args[1] == flag {
			return true
		}
	}
	return false
}

func WalkPaths() ([]string, error) {
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

func GetFilteredResults(currentWorkingDirectory string, input string, dirs []string) []string {
	filteredResults := []string{}
	for _, dir := range dirs {
		if strings.Contains(dir, currentWorkingDirectory) && strings.Contains(strings.ToLower(dir), strings.ToLower(input)) {
			filteredResults = append(filteredResults, dir)
		}
	}
	// matches := fuzzy.Find(input, dirs)
	// filteredResults = append(filteredResults, matches...)
	return filteredResults
}

func OpenInVSCodeFromFinder(selectedItem string, resultlistCount int) {
	if resultlistCount > 0 {
		currentItemName := strings.Split(selectedItem, "/")[len(strings.Split(selectedItem, "/"))-1]
		fmt.Println("To open the directory type:")
		fmt.Println()
		fmt.Print("cd ", selectedItem, "\n")
		fmt.Println()
		fmt.Println("Opening:", currentItemName)
		cmd := exec.Command("code", selectedItem)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		panic("No results found")
	}
}
