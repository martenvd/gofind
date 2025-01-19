// GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o gofind main.go
// GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o gofind main.go

package main

import (
	"flag"
	"log"
	"os"

	"github.com/martenvd/gofind/internal/core"
	"github.com/martenvd/gofind/internal/utils"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	var filteredPath string

	config := utils.CheckConfig(homeDir)
	if config["path"] != nil {
		filteredPath = config["path"].(string)
	}

	updateCache := flag.Bool("u", false, "Whether or not to update the gofind cache.")
	updateCacheFullName := flag.Bool("update", false, "Whether or not to update the gofind cache.")
	flag.Parse()

	if *updateCache || *updateCacheFullName || !utils.FileExists(homeDir+"/.gofind/dirs.txt") {
		cache := core.CheckCache(homeDir)
		dirs, err := utils.WalkPaths(filteredPath, cache)
		if err != nil {
			log.Fatal(err)
		}
		core.CacheDirs(homeDir, dirs)
	}

	dirs, err := core.ReadDirs(homeDir)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		core.Prompt(dirs)
	} else {
		core.Find(dirs)
	}
}
