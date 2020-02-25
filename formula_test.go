package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_formula_String(t *testing.T) {
	f := formula{
		Name:     "test",
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
	if f.isInstalled() {
		t.Fatalf(f.Name)
	}
	f.Name = "hugo"
	if true == f.isInstalled() {
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
