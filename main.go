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
var templateFiles embed.FS
var templates = template.Must(template.ParseFS(templateFiles, "templates/*"))

//go:embed static/*
var staticFiles embed.FS

//go:embed data.csv
var defaultRawDataAsString string

// prevent err is undefined errors
var err error

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

func marshalStrings(data []string) string {
	tmp, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	return string(tmp)
}

func buildIndexPage(rawData [][]string, workoutsJSON string, instructorsJSON string, locationsJSON string) []byte {
	tableHeaders := [5]string{"Date", "Time", "Workout", "Instructor", "Location"}

	classes, err := json.Marshal(rawData[1:])
	if err != nil {
		log.Fatal(err)
	}
	classesJSON := string(classes)

	tmp, err := staticFiles.ReadFile("static/index.css")
	if err != nil {
		log.Fatal()
	}
	styleFile := string(tmp)

	tmp, err = staticFiles.ReadFile("static/index.js")
	if err != nil {
		log.Fatal()
	}
	scriptFile := string(tmp)

	data := struct {
		TableHeaders [5]string
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
	instructorJSON      string
	locationsJSON       string
)

func update() {
	rawData := fetchRawData()
	tmp := extractUniqueEntries(rawData)
	workouts, locations, instructors := marshalStrings(tmp[0]), marshalStrings(tmp[1]), marshalStrings(tmp[2])
	cachedIndexTemplate = buildIndexPage(rawData, workouts, instructors, locations)
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
			log.Fatal()
		}
		stateToLocationsJSON := string(tmp)

		tmp, err = staticFiles.ReadFile("static/index.css")
		if err != nil {
			log.Fatal()
		}
		style := string(tmp)

		tmp, err = staticFiles.ReadFile("static/settings.css")
		if err != nil {
			log.Fatal()
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
		"ALBION PARK", "ANU AQUATICS", "BLACKTOWN", "BURWOOD", "FIVE DOCK", "GLADESVILLE",
		"GOULBURN - LANSDOWNE STREET", "GREGORY HILLS (LIVE WELL)", "HIIT REPUBLIC CORRIMAL",
		"HIIT REPUBLIC QUEANBEYAN", "HIIT REPUBLIC SHELLHARBOUR", "HIIT REPUBLIC WOLLONGONG",
		"LAVINGTON", "MOONEE BEACH", "PENRITH - MULGOA ROAD", "PYRMONT", "QUEANBEYAN", "RHODES",
		"ROSEBERY", "ROUSE HILL", "SHELLHARBOUR", "ST PETERS", "TOORMINA", "WAGGA WAGGA",
		"WOLLONGONG",
	},
	"ACT": {
		"ANU", "BELCONNEN (CISAC PLATINUM)", "BELCONNEN (OATLEY COURT)", "CISAC LADIES ONLY",
		"CONDER", "GOLD CREEK COUNTRY CLUB", "GUNGAHLIN PLATINUM", "HIIT REPUBLIC BRADDON",
		"HIIT REPUBLIC CANBERRA CITY", "HIIT REPUBLIC CISAC", "HIIT REPUBLIC ERINDALE",
		"HIIT REPUBLIC GOULBURN", "HIIT REPUBLIC GUNGAHLIN", "HIIT REPUBLIC KINGSTON",
		"HIIT REPUBLIC MITCHELL", "HIIT REPUBLIC TUGGERANONG", "HIIT REPUBLIC WESTON",
		"HIIT REPUBLIC WODEN", "KAMBAH", "KINGSTON", "KIPPAX", "MAWSON", "MITCHELL", "QUEANBEYAN", "TUGGERANONG",
	},
	"QLD": {
		"ASPLEY", "BROADBEACH", "CLEVELAND", "DEAGON", "HIIT REPUBLIC REDCLIFFE",
		"HIIT REPUBLIC YAMANTO", "IPSWICH", "MOLOOLABA", "NOOSAVILLE", "NORMAN PARK (ACTIVE LIFE)",
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
