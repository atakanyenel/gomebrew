package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type version struct {
	Stable string
}
type fileUrl struct {
	URL    string
	Sha256 string
}
type formula struct {
	Name         string
	Desc         string
	Homepage     string
	Revision     int
	Dependencies []string
	Versions     version
	Bottle       struct {
		Stable struct {
			Files struct {
				Catalina    fileUrl
				Mojave      fileUrl
				High_Sierra fileUrl
			}
		}
	}
}

func (f formula) download() string {
	file := f.getMacOSVersion()
	log.Printf("Downloading from %s", file.URL)
	fp, _ := downloadFile(file.URL)

	if err := checkSha256(fp, file.Sha256); err != nil {
		log.Fatalf("%s", err)
	}
	return fp
}
func (f formula) install(tarLocation string) error {

	check(untarFile(tarLocation))
	if _, err := os.Stat(packagesDir + "/" + f.Name); os.IsNotExist(err) { //check untared file exists
		return err
	}

	//delete tar
	os.Remove(tarLocation)

	// now we only have uncompressed file
	return f.createSymLink()
}

func (f formula) String() string {
	return fmt.Sprintf("%s -> %s\nVersion: %s\nHomepage: %s", f.Name, f.Desc, f.Versions.Stable, f.Homepage)
}

func (f formula) getMacOSVersion() fileUrl {

	out, _ := exec.Command("sw_vers", "-productVersion").Output()
	version := string(out)[:5]
	log.Printf("OS version is %s", version)
	switch version {
	case "10.15":
		return f.Bottle.Stable.Files.Catalina
	case "10.14":
		return f.Bottle.Stable.Files.Mojave
	case "10.13":
		return f.Bottle.Stable.Files.High_Sierra
	}
	return fileUrl{}
}

func (f formula) isInstalled() bool {

	if _, err := exec.LookPath("gome-" + f.Name); err != nil {
		return false
	}
	if _, err := os.Stat(packagesDir + "/" + f.Name); os.IsNotExist(err) { //check untared file exists
		return false
	}
	return true
}

func (f formula) createSymLink() error {
	execPath := fmt.Sprintf("%s/%s/%s/bin/%s", packagesDir, f.Name, f.Versions.Stable, f.Name)
	if _, err := os.Stat(execPath); os.IsNotExist(err) { //check executable exists
		return err
	}

	destination := fmt.Sprintf(GomeSymPath, f.Name)
	if _, err := os.Lstat(destination); err == nil { //if symlink exists, gives error
		os.Remove(destination)
	}
	log.Printf("Creating symlink: %s --> %s", execPath, destination)
	err := os.Symlink(execPath, destination)
	return err
}

func (f formula) uninstall() {
	log.Printf("uninstall called with %s", f.Name)
	destination := fmt.Sprintf(GomeSymPath, f.Name)
	if _, err := os.Lstat(destination); err == nil {
		os.Remove(destination) //remove symlink
	}
	packagePath := fmt.Sprintf("%s/%s", packagesDir, f.Name)
	err := os.RemoveAll(packagePath) //remove gome_package folder
	check(err)
}