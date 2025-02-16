package main

// import (
// 	"log"
// 	"os"
// 	"testing"
// )

// // build template and write to file
// func TestBuildIndexPage(t *testing.T) {
// 	rawData := fetchRawData()
// 	tmp := extractUniqueEntries(rawData)
// 	workouts, locations, instructors := tmp[0], tmp[1], tmp[2]
// 	indexPage := buildIndexPage(rawData, workouts, instructors, locations)

// 	file, err := os.Create("template.html")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()
// 	_, err = file.Write(indexPage)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
