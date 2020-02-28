package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const (
	DeleteResources = iota
	SymlinkResources
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

	f.handleSymlinks(SymlinkResources)
	// now we only have uncompressed file

	return nil
}

func (f formula) String() string {
	output := `
---
{{.Name}}
Desc: {{.Desc}}
Version: {{.Versions.Stable}}
Homepage: {{.Homepage}}
IsInstalled: {{ if .IsInstalled }}✅{{else}}❌{{end}}
---
`
	tmpl, err := template.New("info").Parse(output)
	check(err)
	var tpl bytes.Buffer
	check(tmpl.Execute(&tpl, f))

	return tpl.String()
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

func (f formula) IsInstalled() bool {

	if _, err := exec.LookPath("gome-" + f.Name); err != nil {
		return false
	}
	if _, err := os.Stat(packagesDir + "/" + f.Name); os.IsNotExist(err) { //check untared file exists
		return false
	}
	return true
}

func (f formula) uninstall() {
	log.Printf("uninstall called with %s", f.Name)

	f.handleSymlinks(DeleteResources)

	packagePath := fmt.Sprintf("%s/%s", packagesDir, f.Name)
	err := os.RemoveAll(packagePath) //remove gome_package folder
	check(err)
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

func (f formula) handleSymlinks(action int) error {

	execPath := fmt.Sprintf("%s/%s/%s/bin/%s", packagesDir, f.Name, f.getRealLocation(), f.Name)
	destination := fmt.Sprintf(GomeSymPath, f.Name)

	type resource struct {
		realLocation    string
		symLinkLocation string
	}
	//add executable to resources
	shareResources := []resource{{execPath, destination}}

	r := fmt.Sprintf("%s/%s/%s/share/", packagesDir, f.Name, f.getRealLocation())

	if _, err := os.Stat(r); os.IsNotExist(err) { //check file exists
		return err
	}

	filepath.Walk(r, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			rel, _ := filepath.Rel(r, path)

			g := filepath.Join("/usr/local/share/", rel)

			updatedLocation := filepath.Join(filepath.Dir(g), fmt.Sprintf("gome-%s", filepath.Base(g)))
			shareResources = append(shareResources, resource{path, updatedLocation})

		}
		return nil
	})

	for _, r := range shareResources {
		handleSymLink(r.realLocation, r.symLinkLocation, action)
	}

	return nil
}

func (f formula) getRealLocation() string {
	if f.Linked_keg != "" {
		return f.Linked_keg
	}
	return f.Versions.Stable
}
