package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const (
	DeleteResources = iota
	CreateResources
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
				Big_Sur     fileUrl
			}
		}
	}
	Installed []struct {
		Runtime_dependencies []struct {
			Full_name string
			Version   string
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

	f.handleSymlinks(CreateResources)
	// now we only have uncompressed file

	return nil
}

func (f formula) String() string {
	const output = `
---
{{.Name}}
Desc: {{.Desc}}
Version: {{.Versions.Stable}}
Homepage: {{.Homepage}}
IsInstalled: {{ if .IsInstalled }}✅{{else}}❌{{end}}
---
`
	tmpl, _ := template.New("info").Parse(output)
	var tpl bytes.Buffer
	check(tmpl.Execute(&tpl, f))

	return tpl.String()
}

func (f formula) getMacOSVersion() fileUrl {

	out, _ := exec.Command("sw_vers", "-productVersion").Output()
	version := string(out)[:5]
	log.Printf("OS version is %s", version)
	switch version {
	case "11.0.":
		return f.Bottle.Stable.Files.Big_Sur
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

	packagePath := filepath.Join(packagesDir, f.Name)
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

	type resource struct {
		realLocation    string
		symLinkLocation string
	}
	//add executable to resources
	packageResources := []resource{}

	binPath := filepath.Join(packagesDir, f.Name, f.getRealLocation(), "bin")
	if executables, err := ioutil.ReadDir(binPath); err == nil {
		for _, exe := range executables {

			destination := filepath.Join("/usr/local/bin", "gome-"+exe.Name())
			packageResources = append(packageResources, resource{filepath.Join(binPath, exe.Name()), destination})
		}
	}

	sharePath := filepath.Join(packagesDir, f.Name, f.getRealLocation(), "share")

	if _, err := os.Stat(sharePath); os.IsNotExist(err) { //check file exists
		return err
	}

	filepath.Walk(sharePath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			rel, _ := filepath.Rel(sharePath, path)

			g := filepath.Join("/usr/local/share/", rel)

			updatedLocation := filepath.Join(filepath.Dir(g), "gome-"+filepath.Base(g)) //add gome- to pages
			packageResources = append(packageResources, resource{path, updatedLocation})

		}
		return nil
	})
	for _, r := range packageResources {
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

func (f formula) hasRuntimeDependencies() bool {
	for _, i := range f.Installed {
		if len(i.Runtime_dependencies) != 0 {
			return true
		}
	}
	return false
}
