package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

type ApodResponse struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Date  string `json:"date"`
	Explanation string `json:"explanation"`
}

type RoverPhoto struct {
	ImgSrc    string `json:"img_src"`
	EarthDate string `json:"earth_date"`
}

type RoverResponse struct {
	Photos []RoverPhoto `json:"photos"`
}

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		h(w, r)
	}
}

func apodHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("NASA_API_KEY")
	apodURL := "https://api.nasa.gov/planetary/apod?api_key=" + apiKey
	resp, err := http.Get(apodURL)
	if err != nil {
		http.Error(w, "Error fetching APOD", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error from APOD API", resp.StatusCode)
		return
	}

	var apod ApodResponse
	if err := json.NewDecoder(resp.Body).Decode(&apod); err != nil {
		http.Error(w, "Error decoding APOD response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apod)
}

func roverHandler(w http.ResponseWriter, r *http.Request) {
	sol := r.URL.Query().Get("sol")
	camera := r.URL.Query().Get("camera")
	if sol == "" || camera == "" {
		http.Error(w, "Missing 'sol' or 'camera'", http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(sol); err != nil {
		http.Error(w, "'sol' must be a number", http.StatusBadRequest)
		return
	}
	apiKey := os.Getenv("NASA_API_KEY")
	roverURL := "https://api.nasa.gov/mars-photos/api/v1/rovers/curiosity/photos?sol=" + sol + "&camera=" + camera + "&api_key=" + apiKey

	resp, err := http.Get(roverURL)
	if err != nil {
		http.Error(w, "Error fetching rover photos", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error from Rover API", resp.StatusCode)
		return
	}

	var roverResp RoverResponse
	if err := json.NewDecoder(resp.Body).Decode(&roverResp); err != nil {
		http.Error(w, "Error decoding rover response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roverResp)
}

func main() {
	http.HandleFunc("/api/apod", withCORS(apodHandler))
	http.HandleFunc("/api/rover", withCORS(roverHandler))

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("API server listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}