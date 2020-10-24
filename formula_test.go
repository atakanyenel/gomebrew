package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	teardown()
	os.Exit(code)
}

func teardown() {
	os.RemoveAll("upx-3.96.catalina.bottle.tar.gz")
	os.RemoveAll("upx/")
}
func Test_formula_String(t *testing.T) {
	f := formula{
		Name:     "tree",
		Desc:     "a test package",
		Homepage: "https://test.sh",
		Versions: version{
			Stable: "1.1.1",
		},
	}
	fmt.Println(f)

}

func Test_checkSha256(t *testing.T) {
	fileName := "hello.txt"
	d1 := []byte("Hello World\n")
	err := ioutil.WriteFile(fileName, d1, 0644)
	check(err)
	if nil != checkSha256(fileName, "d2a84f4b8b650937ec8f73cd8be2c74add5a911ba64df27458ed8229da804a26") {
		t.Fail()
	}
	os.Remove(fileName)
}

func Test_downloadFile(t *testing.T) {
	fileURL := "https://homebrew.bintray.com/bottles/upx-3.96.catalina.bottle.tar.gz"

	fp, _ := downloadFile(fileURL)
	if nil != checkSha256(fp, "1089a067bec1387bfa8080565f95601493291933b83a510057ba6f1e7fd06d91") {
		t.Fail()
	}

}

func Test_untarFile(t *testing.T) {

	packagesDir = "."
	filepath := "upx-3.96.catalina.bottle.tar.gz"
	if err := untarFile(filepath); err != nil {
		fmt.Println(err)
		t.Fail()
	}

}

func Test_isInstalled(t *testing.T) {
	f := formula{
		Name:     "shouldnt_be_installed",
		Desc:     "Compress/expand executable files",
		Homepage: "https://upx.github.io/",
		Versions: version{
			Stable: "3.96",
		},
	}
	if f.IsInstalled() {
		t.Fatalf(f.Name)
	}
	f.Name = "hugo"
	if true == f.IsInstalled() {
		t.Fatalf(f.Name)
	}

}

func Test_glob(t *testing.T) {
	files, _ := filepath.Glob("gome_packages/*/*")
	for _, f := range files {
		program, v := filepath.Split(f)
		program = filepath.Base(program)

		f := formula{
			Name:     program,
			Versions: version{Stable: v},
		}

		fmt.Println(f)
	}
}

func Test_betterManPages(t *testing.T) {
	packagesDir, _ = filepath.Abs("gome_packages")

	f := formula{Name: "wget",
		Linked_keg: "1.20.3_2",
		Versions: version{
			Stable: "1.20.3",
		},
	}
	f.handleSymlinks(DeleteResources)

}

func Test_runtimeDeps(t *testing.T) {
	f, _ := getUpstreamFormula("wget")

	if !f.hasRuntimeDependencies() {
		fmt.Println(f.Name)
		t.Fail()
	}

	f, _ = getUpstreamFormula("minikube")

	if f.hasRuntimeDependencies() {
		fmt.Println(f.Name)
		t.Fail()
	}
	fmt.Printf("%+v", f.Installed)
}
func Test_install_deps(t *testing.T) {
	file, _ := ioutil.ReadFile("scripts/formula.json")
	data := []formula{}
	if err := json.Unmarshal([]byte(file), &data); err != nil {
		fmt.Println(err)
	}

	for _, k := range data {
		if len(k.Dependencies) == 1 && !k.hasRuntimeDependencies() {
			if k.Name == "minikube" {
				fmt.Println(k.Name)
			}
			fmt.Println(k.Name)
		}
	}

}

func Test_recursiveInstall(t *testing.T) {
	t_install("minikube")
}
func t_install(program string) {

	g, _ := getUpstreamFormula(program)
	fmt.Println(g.Name)
	for _, k := range g.Dependencies {
		t_install(k)
	}
}
