package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"strings"
	"strconv"
	"time"
	"log"
	"os"
	"github.com/marni/goigc"		// Main library for working on IGC files
	"github.com/p3lim/iso8601"	// For formatting time into ISO 8601
)


var TrackUrl map[int]string		// Declare map for storing URLs
var TrackIds map[int]int			// Declare map for storing IDs corresponding to URL
var startTime time.Time				// Variable for calculating uptime


type TrackInfo struct {				// Encoding complete IGC profile
	H_date string					`json:"date"`
	Pilot string					`json:"pilot"`
	Glider string					`json:"glider"`
	Glider_id string 			`json:"glider_id"`
	Track_length float64	`json:"track_length"`
}


type MetaData struct {		// Encoding meta information of server
	Uptime 	string	`json:"uptime"`
	Info 		string	`json:"info"`
	Version string	`json:"version"`
}


type TrackId struct {			// Encoding URL ID
	Id int	`json: "id"`
}


func handleIgcPlus(w http.ResponseWriter, r *http.Request) {

	if (r.Method == http.MethodGet) {	// Check if GET was called

		parts := strings.Split(r.URL.Path, "/")	// Storing URL parts in a new array
		l := len(parts)				// Number of parts
		t := len(TrackIds)		// Number of already stored IDs
		idString := parts[4]	// Stores ID from 5th element in a string

		id := strconv.Atoi(idString)	// Converts to int

		url := TrackUrl[id]	// Gets the correct URL based on input id

		track, err2 := igc.ParseLocation(url)	// Parses information from URL location
		if(err2 != nil){
				http.Error(w, "Parsing Failed. Request Timeout.", 408)
		}																			// Calculates track length
		trackLen := track.Points[0].Distance(track.Points[len(track.Points)-1])

		// WITH FIELD
		if (l == 6 && id > 0 && id <= t) {					// If the call consists of 6 parts

			field := parts[5]				// Store the field input

			if (field == "pilot") {	// Check which variable 'field' is equal to
				fmt.Fprintln(w, track.Pilot)

			} else if (field == "glider") {
				fmt.Fprintln(w, track.GliderType)

			} else if (field == "glider_id") {
				fmt.Fprintln(w, track.GliderID)

			} else if (field == "track_length") {
				fmt.Fprintln(w, trackLen)

			} else if (field == "H_date") {
				fmt.Fprintln(w, track.Date.String())

			} else {	// If neither of the correct variables are called: 400
				http.Error(w, "No valid Field value.", 400)
			}

		// ONLY ID, EMPTY FIELD
		} else if (l == 5 && id > 0 && id <= t) {	// 5 parts and id input is valid (1-len)
								// Creating temporary struct to hold variables
			temp := TrackInfo{track.Date.String(), track.Pilot, track.GliderType, track.GliderID, trackLen}
								// Encodes temporary struct and shows information on screen
			http.Header.Add(w.Header(), "Content-type", "application/json")
			err3 := json.NewEncoder(w).Encode(temp)
			if(err3 != nil){
			  	http.Error(w, "Encoding Failed. Request Timeout.", 408)
			}
		// If neither ID or ID + Field was found: 400
		} else {
			http.Error(w, "No valid ID value.", 400)	// Bad request
		}
	}
}


func handleIgc(w http.ResponseWriter, r *http.Request) {

		if (r.Method == http.MethodGet) {		// Check if GET was called

			var a []int												// Initialize empty int slice
			a = make([]int, len(TrackUrl))

			for key, url := range TrackUrl {	// Append each key in TrackUrl to the slice 'a'
				a = append(a, key)
				fmt.Println("\n", url)					// To avoid console error of URL not used.
			}
			w.Header().Set("Content-Type", "application/json")

			err := json.NewEncoder(w).Encode(a)
			if(err != nil){
			  	http.Error(w, "Encoding Failed. Request Timeout.", 408)
			}

		} else if (r.Method == http.MethodPost) {	// Check if POST was called

			var temp map[string]interface{}	// Interface is unknown type / C++ auto

			err2 := json.NewDecoder(r.Body).Decode(&temp)	// Decode posted url
			if (err2 != nil) {
				http.Error(w, "Decoding Failed. Request Timeout.", 408)
			}
														// Internal identification: Int up from 1
			tempLen := len(TrackUrl) + 1
														// Places url in map spot nr 1 and up
			TrackUrl[tempLen] = temp["url"].(string)
			TrackIds[tempLen] = tempLen

			idStruct := TrackId{Id: tempLen}
														// Define header for correct output
			http.Header.Add(w.Header(), "Content-type", "application/json")
			err3 := json.NewEncoder(w).Encode(idStruct)
			if (err3 != nil) {
				http.Error(w, "Encoding Failed. Request Timeout.", 408)
			}
		}
}


func handleApi(w http.ResponseWriter, r *http.Request) {
	if (r.Method == http.MethodGet) {	// Check if GET was called
									// Using included library, format time Since into ISO 8601
		t := iso8601.Format(time.Since(startTime))

									// Prepare struct for encoding
		temp := MetaData{t, "Service for IGC tracks.", "v1"}

									// Define header for correct output
		http.Header.Add(w.Header(), "Content-type", "application/json")
		err := json.NewEncoder(w).Encode(temp)
		if(err != nil){
				http.Error(w, "Encoding Failed. Request Timeout.", 408)
		}
	}
}

				// Handle every call beginning with anything else than igcinfo/api
func handleInvalid(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

				// Find the correct Port to run app on Heroku
func GetPort() string {
	 	var port = os.Getenv("PORT")
 				// Port sets to :8080 as a default
 		if (port == "") {
 			port = "8080"
			fmt.Println("No PORT variable detected, defaulting to " + port)
 		}
 		return (":" + port)
}


func main() {

		TrackUrl = make(map[int]string) // Initializing map arrays
		TrackIds = make(map[int]int)
		startTime = time.Now()	// Initializes timer

		http.HandleFunc("/", handleInvalid)
		http.HandleFunc("/igcinfo/api", handleApi)
		http.HandleFunc("/igcinfo/api/igc", handleIgc)
		http.HandleFunc("/igcinfo/api/igc/", handleIgcPlus)

		err := http.ListenAndServe(GetPort(), nil)
		if err != nil {
				log.Fatal("ListenAndServe Error: ", err)
		}
}
