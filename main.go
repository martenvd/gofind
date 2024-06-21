// GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o gofind main.go
package main

import (
	"flag"
	"log"
	"os"

	"github.com/martenvd/gofind/internal/app"
	"github.com/martenvd/gofind/internal/utils"
)

func main() {

	updateCache := flag.Bool("u", false, "Whether or not to update the gofind cache.")
	updateCacheFullName := flag.Bool("update", false, "Whether or not to update the gofind cache.")
	flag.Parse()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	if *updateCache || *updateCacheFullName || !utils.FileExists(homeDir+"/.gofind/dirs.txt") {
		dirs, err := utils.WalkPaths()
		if err != nil {
			log.Fatal(err)
		}
		app.CacheDirs(homeDir, dirs)
	}

	dirs, err := app.ReadDirs(homeDir)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(dirs)
	// prompt(dirs)
	app.FuzzyFind(dirs)
}
