package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"text/template"
	"time"
)

var (
	rawData     [][]string
	lastFetched time.Time
)
var (
	instructorsJSON string
	locationsJSON   string
	workoutsJSON    string
)

func fetchCSV() {
	file, err := os.Open("data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	rawData = data
	lastFetched = time.Now()
}

func jsonifySet(set map[string]struct{}) string {
	uniqueArray := []string{}
	for value := range set {
		uniqueArray = append(uniqueArray, value)
	}
	sort.Strings(uniqueArray)
	json, err := json.Marshal(uniqueArray)
	if err != nil {
		log.Fatal(err)
	}
	return string(json)
}

// update the data based on the current rawData
func updateData() {
	if len(rawData) == 0 {
		return
	}
	locations := map[string]struct{}{}
	instructors := map[string]struct{}{}
	workouts := map[string]struct{}{}
	for _, row := range rawData[1:] {
		workout, instructor, location := row[2], row[3], row[4]
		workouts[workout] = struct{}{}
		instructors[instructor] = struct{}{}
		locations[location] = struct{}{}
	}
	workoutsJSON = jsonifySet(workouts)
	instructorsJSON = jsonifySet(instructors)
	locationsJSON = jsonifySet(locations)
}

func main() {
	if lastFetched.IsZero() || time.Since(lastFetched) > 24*time.Hour {
		fetchCSV()
		updateData()
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Fatal(err)
		}

		tableHeaders := [5]string{"Date", "Time", "Workout", "Instructor", "Location"}

		classes, err := json.Marshal(rawData[1:])
		if err != nil {
			log.Fatal(err)
		}
		classesJSON := string(classes)

		data := struct {
			TableHeaders [5]string
			Classes      string
			Locations    string
			Workouts     string
		}{
			TableHeaders: tableHeaders,
			Classes:      classesJSON,
			Locations:    locationsJSON,
			Workouts:     workoutsJSON,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
	})

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// var stateToClubs = map[string][]string{
// 	"NSW": {
// 		"ALBION PARK", "ANU AQUATICS", "BLACKTOWN", "BURWOOD", "FIVE DOCK", "GLADESVILLE",
// 		"GOULBURN - LANSDOWNE STREET", "GREGORY HILLS (LIVE WELL)", "HIIT REPUBLIC CORRIMAL",
// 		"HIIT REPUBLIC QUEANBEYAN", "HIIT REPUBLIC SHELLHARBOUR", "HIIT REPUBLIC WOLLONGONG",
// 		"LAVINGTON", "MOONEE BEACH", "PENRITH - MULGOA ROAD", "PYRMONT", "RHODES",
// 		"ROSEBERY", "ROUSE HILL", "SHELLHARBOUR", "ST PETERS", "TOORMINA", "WAGGA WAGGA",
// 		"WOLLONGONG",
// 	},
// 	"ACT": {
// 		"ANU", "BELCONNEN (CISAC PLATINUM)", "BELCONNEN (OATLEY COURT)", "CISAC LADIES ONLY",
// 		"CONDER", "GOLD CREEK COUNTRY CLUB", "GUNGAHLIN PLATINUM", "HIIT REPUBLIC BRADDON",
// 		"HIIT REPUBLIC CANBERRA CITY", "HIIT REPUBLIC CISAC", "HIIT REPUBLIC ERINDALE",
// 		"HIIT REPUBLIC GOULBURN", "HIIT REPUBLIC GUNGAHLIN", "HIIT REPUBLIC KINGSTON",
// 		"HIIT REPUBLIC MITCHELL", "HIIT REPUBLIC TUGGERANONG", "HIIT REPUBLIC WESTON",
// 		"HIIT REPUBLIC WODEN", "KAMBAH", "KINGSTON", "KIPPAX", "MAWSON", "MITCHELL", "TUGGERANONG",
// 		"QUEANBEYAN",
// 	},
// 	"QLD": {
// 		"ASPLEY", "BROADBEACH", "CLEVELAND", "DEAGON", "HIIT REPUBLIC REDCLIFFE",
// 		"HIIT REPUBLIC YAMANTO", "IPSWICH", "MOLOOLABA", "NOOSAVILLE", "NORMAN PARK (ACTIVE LIFE)",
// 		"REDCLIFFE", "SIPPY DOWNS", "SUNNYBANK HILLS (HEALTHWORKS)", "TENERIFFE", "WEST END",
// 	},
// 	"WA": {
// 		"BUTLER", "CLARKSON", "KINGS SQUARE", "MANDURAH", "PORT KENNEDY", "WANNEROO",
// 	},
// 	"VIC": {
// 		"CARIBBEAN PARK", "HIIT REPUBLIC COBURG", "MALVERN", "MENTONE", "MULGRAVE", "OAKLEIGH",
// 		"PARKDALE", "SCORESBY", "SHEPPARTON", "SOUTH MORANG (ONE HEALTH)", "UPWEY", "WILLIAMSTOWN",
// 		"WODONGA",
// 	},
// 	"NT": {
// 		"COOLALINGA (IFITNESS 247)", "DARWIN CITY (IFITNESS 247)", "MILLNER (IFITNESS 247)",
// 		"PALMERSTON (IFITNESS 247)",
// 	},
// }
