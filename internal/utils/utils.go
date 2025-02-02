package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
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

func WalkPaths(filteredPath string, cache []string) ([]string, error) {
	var dirs []string
	var mu sync.Mutex
	var wg sync.WaitGroup

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

	if filteredPath != "" {
		currentDir = filteredPath
	}

	filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		for _, dir := range cache {
			if path == dir {
				dirs = append(dirs, path)
				return nil
			}
		}

		if err != nil {
			fmt.Println(err)
		}

		if info.IsDir() {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				files, err := os.ReadDir(path)
				if err != nil {
					fmt.Println(err)
				}

				for _, file := range files {
					if rgx.MatchString(file.Name()) {
						mu.Lock()
						dirs = append(dirs, path)
						mu.Unlock()
						break
					}
				}
			}(path)
		}

		return nil
	})

	wg.Wait()

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

func CheckConfig(homeDir string) map[string]interface{} {
	file, err := os.Open(homeDir + "/.gofind/config.json")
	if err != nil {
		return nil
	}

	var config map[string]interface{}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}

	if config["path"] == nil {
		fmt.Println("No version found in config")
		return nil
	}

	defer file.Close()

	return config
}
