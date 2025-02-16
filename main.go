package main

import (
	"bytes"
	"embed"
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

//go:embed templates/*
var resources embed.FS
var t = template.Must(template.ParseFS(resources, "templates/*"))

//go:embed data.csv
var defaultRawDataAsString string

func fetchRawData() [][]string {
	reader := csv.NewReader(strings.NewReader(defaultRawDataAsString))
	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return data
}

// returns [workouts, locations, instructors],
// 3 arrays of sorted unique values
func extractUniqueEntries(rawData [][]string) [3][]string {
	locationsSet, instructorsSet, workoutsSet := map[string]struct{}{}, map[string]struct{}{}, map[string]struct{}{}
	for _, row := range rawData[1:] {
		workout, instructor, location := row[2], row[3], row[4]
		workoutsSet[workout] = struct{}{}
		instructorsSet[instructor] = struct{}{}
		locationsSet[location] = struct{}{}
	}

	setToSortedArray := func(set map[string]struct{}) []string {
		uniqueArray := []string{}
		for value := range set {
			uniqueArray = append(uniqueArray, value)
		}
		// sort.Strings(uniqueArray)
		return uniqueArray
	}

	workouts := setToSortedArray(workoutsSet)
	locations := setToSortedArray(locationsSet)
	instructors := setToSortedArray(instructorsSet)
	return [3][]string{workouts, locations, instructors}
}

func jsonify(data []string) string {
	tmp, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	return string(tmp)
}

func buildIndexPage(rawData [][]string, workouts []string, instructors []string, locations []string) []byte {
	tableHeaders := [5]string{"Date", "Time", "Workout", "Instructor", "Location"}

	classes, err := json.Marshal(rawData[1:])
	if err != nil {
		log.Fatal(err)
	}
	classesJSON := string(classes)

	workoutsJSON := jsonify(workouts)
	instructorsJSON := jsonify(instructors)
	locationsJSON := jsonify(locations)

	data := struct {
		TableHeaders [5]string
		Classes      string
		Locations    string
		Workouts     string
		Instructors  string
	}{
		TableHeaders: tableHeaders,
		Classes:      classesJSON,
		Locations:    locationsJSON,
		Workouts:     workoutsJSON,
		Instructors:  instructorsJSON,
	}

	var buf bytes.Buffer

	err = t.ExecuteTemplate(&buf, "index.html", data)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

var (
	cachedIndexTemplate []byte
)

func update() {
	rawData := fetchRawData()
	tmp := extractUniqueEntries(rawData)
	workouts, locations, instructors := tmp[0], tmp[1], tmp[2]
	cachedIndexTemplate = buildIndexPage(rawData, workouts, instructors, locations)
}

func main() {
	println("main running")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// if lastFetched.IsZero() || time.Since(lastFetched) > 24*time.Hour {
	// 	fetchRawData()
	// 	updateData()
	// }
	update()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(cachedIndexTemplate)
	})

	log.Printf("Server started on :%s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
