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

	f, err := getFormula(program)
	check(err)

	if len(f.Dependencies) != 0 || f.Revision != 0 {
		log.Println("Gomebrew currently does not support dependencies or packages with revisions")
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
	if f, err := getFormula(program); err == nil {
		fmt.Print(f)
	}
}

func uninstall(program string) {
	formula{Name: program}.uninstall()
}

func upgrade() {
	formulas := getInstalledFormulas()

	for _, current := range formulas {
		upstream, err := getFormula(current.Name)
		check(err)
		if upstream.Versions.Stable == current.Versions.Stable {
			log.Printf("%s is already the latest version: %s", current.Name, current.Versions.Stable)
		} else {
			tarLocation := upstream.download()
			check(upstream.install(tarLocation))
		}
	}
}

func purge() {
	for _, current := range getInstalledFormulas() {
		current.uninstall()
	}
	os.Remove(packagesDir)
}
