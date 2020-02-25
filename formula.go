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
	Linked_keg   string
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
	execPath := f.getExecutable()
	destination := fmt.Sprintf(GomeSymPath, f.Name)
	f.createSymLink(execPath, destination)
	manPage := f.getManpage()
	destination = fmt.Sprintf(GomeManPageSymPath, f.Name)
	f.createSymLink(manPage, destination)

	return nil
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

func (f formula) createSymLink(localPath, destination string) error {

	if _, err := os.Stat(localPath); os.IsNotExist(err) { //check file exists
		return err
	}

	if _, err := os.Lstat(destination); err == nil { //if symlink exists, gives error
		os.Remove(destination)
	}
	log.Printf("Creating symlink: %s --> %s", localPath, destination)
	err := os.Symlink(localPath, destination)
	return err
}

func (f formula) uninstall() {
	log.Printf("uninstall called with %s", f.Name)
	destination := fmt.Sprintf(GomeSymPath, f.Name)
	if _, err := os.Lstat(destination); err == nil {
		os.Remove(destination) //remove symlink
	}

	destination = fmt.Sprintf(GomeManPageSymPath, f.Name)
	if _, err := os.Lstat(destination); err == nil {
		os.Remove(destination) //remove man page symlink
	}
	packagePath := fmt.Sprintf("%s/%s", packagesDir, f.Name)
	err := os.RemoveAll(packagePath) //remove gome_package folder
	check(err)
}

func (f formula) getExecutable() string {
	installedVersion := f.Versions.Stable
	if f.Linked_keg != "" {
		installedVersion = f.Linked_keg
	}

	return fmt.Sprintf("%s/%s/%s/bin/%s", packagesDir, f.Name, installedVersion, f.Name)

}

func (f formula) getManpage() string {
	//ln -s  ~/Desktop/Computer_Science/go/src/github.com/atakanyenel/gomebrew/gome_packages/tree/1.8.0/share/man/man1/tree.1  /usr/local/share/man/man1/gome-tree.1
	installedVersion := f.Versions.Stable
	if f.Linked_keg != "" {
		installedVersion = f.Linked_keg
	}

	return fmt.Sprintf("%s/%s/%s/share/man/man1/%s.1", packagesDir, f.Name, installedVersion, f.Name) //too many hardcoded values maybe error here

}

func (f formula) updateExecutable() {
	upstream, err := getUpstreamFormula(f.Name)
	check(err)
	if upstream.Versions.Stable == f.Versions.Stable {
		log.Printf("%s is already the latest version: %s", f.Name, f.Versions.Stable)
	} else {
		tarLocation := upstream.download()
		check(upstream.install(tarLocation))
	}
}
