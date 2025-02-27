package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"testing"
)

// build index page and write it to file
func TestBuildIndexPage(t *testing.T) {
	rawData := readLatest()
	tmp := extractUniqueEntries(rawData)
	workouts, locations, instructors := marshalStrings(tmp[0]), marshalStrings(tmp[1]), marshalStrings(tmp[2])
	indexPage := buildIndexPage(rawData, workouts, instructors, locations)

	file, err := os.Create("template.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	_, err = file.Write(indexPage)
	if err != nil {
		log.Fatal(err)
	}
}

func TestXd(t *testing.T) {
	err := fs.WalkDir(staticFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestXd2(t *testing.T) {
	tmp, err := staticFiles.ReadFile("static/index.js")
	if err != nil {
		log.Fatal()
	}
	scriptFile := string(tmp)
	fmt.Println(scriptFile)
}
