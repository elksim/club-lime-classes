package main

import (
	"bytes"
	"embed"
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"text/template"
	"time"
)

//go:embed templates/*
var templateFiles embed.FS
var templates = template.Must(template.ParseFS(templateFiles, "templates/*"))

//go:embed static/*
var staticFiles embed.FS

// 'hack' to prevent err is undefined errors..
var err error

var rawDataFolder = "data/rawClasses"

// read the latest data from data/
func readLatest() [][]string {
	filePath := getLatestRawDataFilePath()
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(bytes.NewReader(file))
	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return data
}

// return the most recently modified filepath in data/rawClasses/
func getLatestRawDataFilePath() string {
	var mostRecentFile string
	var mostRecentTime time.Time

	err := filepath.WalkDir(rawDataFolder, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if info.ModTime().After(mostRecentTime) {
			mostRecentTime = info.ModTime()
			mostRecentFile = path
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	if mostRecentFile == "" {
		return ""
	}
	return mostRecentFile
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

func marshalStrings(data []string) string {
	sort.Strings(data)
	tmp, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	return string(tmp)
}

func buildIndexPage(rawData [][]string, workoutsJSON string, instructorsJSON string, locationsJSON string) []byte {
	tableHeaders := [4]string{"Time", "Workout", "Instructor", "Location"}

	classes, err := json.Marshal(rawData[1:])
	if err != nil {
		log.Fatal(err)
	}

	classesJSON := string(classes)

	tmp, err := staticFiles.ReadFile("static/index.css")
	if err != nil {
		log.Fatal(err)
	}
	styleFile := string(tmp)

	tmp, err = staticFiles.ReadFile("static/index.js")
	if err != nil {
		log.Fatal(err)
	}
	scriptFile := string(tmp)

	data := struct {
		TableHeaders [4]string
		Classes      string
		Locations    string
		Workouts     string
		Instructors  string
		Script       string
		Style        string
	}{
		TableHeaders: tableHeaders,
		Classes:      classesJSON,
		Locations:    locationsJSON,
		Workouts:     workoutsJSON,
		Instructors:  instructorsJSON,
		Script:       scriptFile,
		Style:        styleFile,
	}

	var buf bytes.Buffer
	err = templates.ExecuteTemplate(&buf, "index.html", data)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

// note: was thinking I also generate and cache the html for the settings page
// if i do that I don't actually need workoutsJSON etal to be global they be scoped to update()
var (
	cachedIndexTemplate []byte
	workoutsJSON        string
	instructorsJSON     string
	locationsJSON       string
)

func update() {
	rawData := readLatest()
	tmp := extractUniqueEntries(rawData)
	workoutsJSON, locationsJSON, instructorsJSON = marshalStrings(tmp[0]), marshalStrings(tmp[1]), marshalStrings(tmp[2])
	cachedIndexTemplate = buildIndexPage(rawData, workoutsJSON, instructorsJSON, locationsJSON)
}

func main() {
	println("main running")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// if lastFetched.IsZero() || time.Since(lastFetched) > 24*time.Hour {
	update()
	// }

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(cachedIndexTemplate)
	})

	http.HandleFunc("/settings/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		tmp, err := json.Marshal(stateToLocations)
		if err != nil {
			log.Fatal(err)
		}
		stateToLocationsJSON := string(tmp)

		tmp, err = staticFiles.ReadFile("static/index.css")
		if err != nil {
			log.Fatal(err)
		}
		style := string(tmp)

		tmp, err = staticFiles.ReadFile("static/settings.css")
		if err != nil {
			log.Fatal(err)
		}
		style2 := string(tmp)

		tmp, err = staticFiles.ReadFile("static/settings.js")
		if err != nil {
			log.Fatal(err)
		}
		script := string(tmp)

		data := struct {
			Workouts         string
			Locations        string
			StateToLocations string
			Script           string
			Style            string
			Style2           string
		}{
			Workouts:         workoutsJSON,
			Locations:        locationsJSON,
			StateToLocations: stateToLocationsJSON,
			Script:           script,
			Style:            style,
			Style2:           style2,
		}
		err = templates.ExecuteTemplate(w, "settings.html", data)
		if err != nil {
			log.Fatal(err)
		}
	})

	log.Printf("Server started on :%s", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Queanbeyan is in both NSW and ACT for pratical reasons - I don't think there are any other
// similar cases to consider.
var stateToLocations = map[string][]string{
	"NSW": {
		"ALBION PARK", "BLACKTOWN", "BURWOOD", "CAMPBELLTOWN", "FIVE DOCK", "GLADESVILLE",
		"GOULBURN - LANSDOWNE STREET", "GREGORY HILLS (LIVE WELL)", "HIIT REPUBLIC CORRIMAL",
		"HIIT REPUBLIC QUEANBEYAN", "HIIT REPUBLIC SHELLHARBOUR", "HIIT REPUBLIC WOLLONGONG",
		"KENNEDY PARK", "LAVINGTON", "MOONEE BEACH", "PENRITH - MULGOA ROAD", "PYRMONT", "QUEANBEYAN", "RHODES",
		"ROSEBERY", "ROUSE HILL", "SHELLHARBOUR", "ST PETERS", "TOORMINA", "WAGGA WAGGA",
		"WOLLONGONG",
	},
	"ACT": {
		"ANU", "ANU AQUATICS", "BELCONNEN (CISAC PLATINUM)", "BELCONNEN (OATLEY COURT)", "CISAC LADIES ONLY",
		"CONDER", "GOLD CREEK COUNTRY CLUB", "GUNGAHLIN PLATINUM", "HIIT REPUBLIC BRADDON",
		"HIIT REPUBLIC CANBERRA CITY", "HIIT REPUBLIC CISAC", "HIIT REPUBLIC ERINDALE",
		"HIIT REPUBLIC GOULBURN", "HIIT REPUBLIC GUNGAHLIN", "HIIT REPUBLIC KINGSTON",
		"HIIT REPUBLIC MITCHELL", "HIIT REPUBLIC TUGGERANONG", "HIIT REPUBLIC WESTON",
		"HIIT REPUBLIC WODEN", "KAMBAH", "KINGSTON", "KIPPAX", "MAWSON", "MITCHELL", "QUEANBEYAN", "TUGGERANONG",
	},
	"QLD": {
		"ASPLEY", "BROADBEACH", "CLEVELAND", "DEAGON", "HIIT REPUBLIC REDCLIFFE",
		"HIIT REPUBLIC YAMANTO", "IPSWICH", "MOOLOOLABA", "NOOSAVILLE", "NORMAN PARK (ACTIVE LIFE)",
		"REDCLIFFE", "SIPPY DOWNS", "SUNNYBANK HILLS (HEALTHWORKS)", "TENERIFFE", "WEST END",
	},
	"WA": {
		"BUTLER", "CLARKSON", "KINGS SQUARE", "MANDURAH", "PORT KENNEDY", "WANNEROO",
	},
	"VIC": {
		"CARIBBEAN PARK", "HIIT REPUBLIC COBURG", "MALVERN", "MENTONE", "MULGRAVE", "OAKLEIGH",
		"PARKDALE", "SCORESBY", "SHEPPARTON", "SOUTH MORANG (ONE HEALTH)", "UPWEY", "WILLIAMSTOWN",
		"WODONGA",
	},
	"NT": {
		"COOLALINGA (IFITNESS 247)", "DARWIN CITY (IFITNESS 247)", "MILLNER (IFITNESS 247)",
		"PALMERSTON (IFITNESS 247)",
	},
}
