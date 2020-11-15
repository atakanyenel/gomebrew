package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var packagesDir string

func main() {

	ex, err := os.Executable()
	check(err)
	packagesDir = filepath.Join(filepath.Dir(ex), "gome_packages")

	_ = os.Mkdir(packagesDir, os.ModePerm) //create folder if not there
	runApp()
}

func runApp() {
	if len(os.Args) < 2 {
		helpCommand()
		return
	}
	command := os.Args[1]
	for _, i := range COMMANDS {
		if command == i.Name {
			i.Action()
			return
		}
	}
	helpCommand()
}

func helpCommand() {
	fmt.Println("Gomebrew - A lite homebrew clone")
	fmt.Println("-------")
	for _, c := range COMMANDS {
		fmt.Printf("%s - %s\n", c.Name, c.Usage)
	}
}

type command struct {
	Name   string
	Usage  string
	Action func()
}

var COMMANDS = []command{
	{
		Name:  "install",
		Usage: "install homebrew package",
		Action: func() {
			for _, p := range os.Args[2:] {
				install(p)
			}
		},
	},
	{
		Name:  "uninstall",
		Usage: "uninstalls homebrew package",
		Action: func() {
			for _, p := range os.Args[2:] {
				uninstall(p)
			}
		},
	},
	{
		Name:  "info",
		Usage: "information about package",
		Action: func() {
			for _, p := range os.Args[2:] {
				info(p)
			}
		},
	},
	{
		Name:  "list",
		Usage: "list installed packages",
		Action: func() {
			list()
		},
	},
	{
		Name:  "prune",
		Usage: "deletes all packages",
		Action: func() {
			prune()
		},
	},
	{
		Name:  "upgrade",
		Usage: "upgrades all packages. If package name given, upgrades only those packages.",
		Action: func() {
			upgrade(os.Args[2:])
		},
	},
}
