package main

import (
	"log"
	"os"
	"path/filepath"
)

const (
	COMMAND_POS = 1
	PACKAGE_POS = 2
	homebrewAPI = "https://formulae.brew.sh/api/formula/%s.json"
	GomeSymPath = "/usr/local/bin/gome-%s"
)

var packagesDir string

func init() { //so that we have packagesDir already defined for tests
	var err error
	packagesDir, err = filepath.Abs("gome_packages")
	check(err)
	_ = os.Mkdir(packagesDir, os.ModePerm) //create folder if not there
}

func main() {
	log.Println("Hello gomebrew")
	commandToFunc := map[string]func(string){"install": install, "info": info, "uninstall": uninstall}

	switch command := os.Args[COMMAND_POS]; command { //todo: add a normal cli
	case "list":
		list()
	case "upgrade":
		upgrade()
	case "purge":
		purge()
	default:
		if fun, ok := commandToFunc[command]; ok {
			programs := os.Args[PACKAGE_POS:]
			for _, p := range programs {

				fun(p)
			}
		} else {
			log.Fatalf("%s is not a command", command)
		}
	}
}
