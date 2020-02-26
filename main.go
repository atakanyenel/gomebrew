package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

const (
	homebrewAPI        = "https://formulae.brew.sh/api/formula/%s.json"
	GomeSymPath        = "/usr/local/bin/gome-%s"
	GomeManPageSymPath = "/usr/local/share/man/man1/gome-%s.1"
)

var packagesDir string

func init() { //so that we have packagesDir already defined for tests
	var err error
	packagesDir, err = filepath.Abs("gome_packages")
	check(err)
	_ = os.Mkdir(packagesDir, os.ModePerm) //create folder if not there
}

func main() {
	app := &cli.App{
		Name:  "gomebrew",
		Usage: "a lite homebrew client",
		Commands: []*cli.Command{
			{
				Name:    "install",
				Usage:   "install homebrew package",
				Aliases: []string{"i"},
				Action: func(c *cli.Context) error {
					for _, p := range c.Args().Slice() {
						install(p)
					}
					return nil
				},
			},
			{
				Name:  "uninstall",
				Usage: "uninstalls homebrew package",
				Action: func(c *cli.Context) error {
					for _, p := range c.Args().Slice() {
						uninstall(p)
					}
					return nil
				},
			},
			{
				Name:  "info",
				Usage: "information about package",
				Action: func(c *cli.Context) error {
					for _, p := range c.Args().Slice() {
						info(p)
					}
					return nil
				},
			},
			{
				Name:    "list",
				Usage:   "list installed packages",
				Aliases: []string{"l"},
				Action: func(c *cli.Context) error {
					list()
					return nil
				},
			},
			{
				Name:  "prune",
				Usage: "deletes all packages",
				Action: func(c *cli.Context) error {
					prune()
					return nil
				},
			},
			{
				Name:  "upgrade",
				Usage: "upgrades all packages. If package name given, upgrades only those packages.",
				Action: func(c *cli.Context) error {
					upgrade(c.Args().Slice())
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
