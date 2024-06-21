package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
