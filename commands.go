package main

import (
	"fmt"
	"log"
	"os"
)

func install(program string) {
	log.Printf("install called with %s", program)

	if (formula{Name: program}.isInstalled()) {
		log.Fatalf("%s is already installed", program)
	}

	f, err := getUpstreamFormula(program)
	check(err)

	if len(f.Dependencies) != 0 {
		log.Println("Gomebrew currently does not support packages with dependencies")
		return // don't fail, so other arguments can continue
	}
	tarLocation := f.download()
	check(f.install(tarLocation))

}

func list() {
	for _, formula := range getInstalledFormulas() {
		fmt.Printf("%s -> %s\n", formula.Name, formula.Versions.Stable)
	}
}

func info(program string) {
	log.Printf("info called with %s", program)
	if f, err := getUpstreamFormula(program); err == nil {
		fmt.Print(f)
	}
}

func uninstall(program string) {
	formula{Name: program}.uninstall()
}

func upgrade(programs []string) {
	formulas := getInstalledFormulas()
	if len(programs) == 0 {
		for _, current := range formulas {
			current.updateExecutable()
		}
	} else {
		for _, current := range programs {
			if f, ok := formulas[current]; ok {
				f.updateExecutable()
			} else {
				log.Printf("%s is not installed", current)
			}
		}
	}
}

func prune() {
	for _, current := range getInstalledFormulas() {
		current.uninstall()
	}
	os.Remove(packagesDir)
}
