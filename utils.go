package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getUpstreamFormula(packageName string) (formula, error) {
	log.Printf("getFormula called with %s", packageName)
	apiURL := fmt.Sprintf(homebrewAPI, packageName)
	log.Printf("Sending request to %s", apiURL)
	r, err := http.Get(apiURL)
	check(err)
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	check(err)
	var f formula
	err = json.Unmarshal(body, &f)
	return f, err
}

func getInstalledFormulas() map[string]formula {
	formulas := map[string]formula{}
	files, _ := filepath.Glob(packagesDir + "/*/*")
	for _, f := range files {
		program, v := filepath.Split(f)
		program = filepath.Base(program)

		installedFormula := formula{
			Name:       program,
			Versions:   version{Stable: strings.Split(v, "_")[0]},
			Linked_keg: v,
		}
		formulas[program] = installedFormula
	}
	return formulas
}

func checkSha256(filepath string, wantedHash string) error {
	contents, err := ioutil.ReadFile(filepath) //Fixme we can also find better function for this
	check(err)
	sum := sha256.Sum256(contents)
	fileChecksum := hex.EncodeToString(sum[:])
	if fileChecksum != wantedHash { // Fixme: we can find a better functions for this

		return fmt.Errorf("Hashsum check failed. Want: %s, Got: %s", wantedHash, fileChecksum)

	}
	log.Println("Hashsum check passed")
	return nil
}

func untarFile(filepath string) error {

	if _, err := exec.LookPath("tar"); err != nil {
		log.Fatal("tar command not found")
	}
	_, err := exec.Command("bash", "-c", fmt.Sprintf("tar xvzf %s -C %s", filepath, packagesDir)).Output()
	return err

}

func handleSymLink(localPath, destination string, action int) error {

	if _, err := os.Lstat(destination); err == nil { //if symlink exists, gives delete it
		os.Remove(destination)
	}

	if action == SymlinkResources {

		if _, err := os.Stat(localPath); os.IsNotExist(err) { //check file exists
			return err
		}
		log.Printf("Creating symlink: %s --> %s", localPath, destination)
		err := os.Symlink(localPath, destination)
		return err
	} //no need for deletesymlink as it's always done first

	return nil
}
