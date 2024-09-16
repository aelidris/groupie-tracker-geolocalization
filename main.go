package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Artist struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Members []string `json:"members"`
}

// type Location struct {
// 	ID        int      `json:"id"`
// 	Locations []string `json:"locations"`
// 	DatesURL  string   `json:"dates"`
// }

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	LatLngs   []LatLng `json:"latlngs"` // To store the coordinates
}

type LatLng struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func getCoordinates(location string) (float64, float64, error) {
	geoapifyURL := fmt.Sprintf("https://api.geoapify.com/v1/geocode/search?text=%s&apiKey=%s", location, "84a9a08245a141e299e5a1fd45b3dbd8")
	resp, err := http.Get(geoapifyURL)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, 0, err
	}

	features := result["features"].([]interface{})
	if len(features) > 0 {
		geometry := features[0].(map[string]interface{})["geometry"].(map[string]interface{})
		coordinates := geometry["coordinates"].([]interface{})
		return coordinates[1].(float64), coordinates[0].(float64), nil
	}

	return 0, 0, fmt.Errorf("no coordinates found")
}

// FetchArtists fetches all artist data from the API
func FetchArtists() ([]Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var artists []Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

// HomePage displays a welcome message and a list of artists
func HomePage(w http.ResponseWriter, r *http.Request) {
	artists, err := FetchArtists()
	if err != nil {
		http.Error(w, "Failed to load artists", http.StatusInternalServerError)
		return
	}

	// Display the homepage with links to artist pages
	fmt.Fprintf(w, "<h1>Welcome to Groupie Trackers</h1>")
	fmt.Fprintf(w, "<h2>Artists List</h2><ul>")
	for _, artist := range artists {
		fmt.Fprintf(w, `<li><a href="/artist?id=%d">%s</a></li>`, artist.ID, artist.Name)
	}
	fmt.Fprintf(w, "</ul>")
}

func ArtistPage(w http.ResponseWriter, r *http.Request) {
	// Extract artist ID from query parameters
	artistID := r.URL.Query().Get("id")
	if artistID == "" {
		http.Error(w, "Artist ID is missing", http.StatusBadRequest)
		return
	}

	// Fetch artist data from the API
	artistURL := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%s", artistID)
	resp, err := http.Get(artistURL)
	if err != nil {
		http.Error(w, "Error fetching artist data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode artist data
	var artist Artist
	err = json.NewDecoder(resp.Body).Decode(&artist)
	if err != nil {
		http.Error(w, "Error decoding artist data", http.StatusInternalServerError)
		return
	}

	// Fetch locations data for the artist
	locationURL := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%s", artistID)
	locationResp, err := http.Get(locationURL)
	if err != nil {
		http.Error(w, "Error fetching location data", http.StatusInternalServerError)
		return
	}
	defer locationResp.Body.Close()

	// Decode location data
	var location Location
	err = json.NewDecoder(locationResp.Body).Decode(&location)
	if err != nil {
		http.Error(w, "Error decoding location data", http.StatusInternalServerError)
		return
	}

	// Fetch coordinates for each location
	var latLngs []LatLng
	for _, loc := range location.Locations {
		lat, lng, err := getCoordinates(loc)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching coordinates for %s: %v", loc, err), http.StatusInternalServerError)
			return
		}
		latLngs = append(latLngs, LatLng{Latitude: lat, Longitude: lng})
	}

	// Use the fetched coordinates to generate a map URL
	mapURL := generateMapURL(latLngs)

	// Display the map and artist information in HTML
	fmt.Fprintf(w, `<html>
		<head><title>%s's Concert Locations</title></head>
		<body>
			<h1>%s</h1>
			<img src="%s" alt="Concert Locations Map">
			<h2>Members:</h2>
			<ul>`, artist.Name, artist.Name, mapURL)
	for _, member := range artist.Members {
		fmt.Fprintf(w, "<li>%s</li>", member)
	}
	fmt.Fprintf(w, `</ul></body></html>`)
}

func generateMapURL(latLngs []LatLng) string {
	baseURL := "https://maps.geoapify.com/v1/staticmap?apiKey=84a9a08245a141e299e5a1fd45b3dbd8&width=800&height=600"
	for _, latlng := range latLngs {
		baseURL += fmt.Sprintf("&marker=lonlat:%f,%f;color:red;size:medium", latlng.Longitude, latlng.Latitude)
	}
	return baseURL
}

// ArtistPage displays information about an artist based on the query parameter 'id'
// func ArtistPage(w http.ResponseWriter, r *http.Request) {
// 	// Extract the 'id' query parameter
// 	idStr := r.URL.Query().Get("id")
// 	if idStr == "" {
// 		http.Error(w, "Artist ID is required", http.StatusBadRequest)
// 		return
// 	}

// 	// Convert the ID to an integer
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Fetch all artists
// 	artists, err := FetchArtists()
// 	if err != nil {
// 		http.Error(w, "Failed to load artists", http.StatusInternalServerError)
// 		return
// 	}

// 	// Find the artist with the matching ID
// 	var foundArtist *Artist
// 	for _, artist := range artists {
// 		if artist.ID == id {
// 			foundArtist = &artist
// 			break
// 		}
// 	}

// 	// If artist is not found
// 	if foundArtist == nil {
// 		http.Error(w, "Artist not found", http.StatusNotFound)
// 		return
// 	}

// 	// Display the artist information
// 	fmt.Fprintf(w, "<h1>Artist: %s</h1>", foundArtist.Name)
// 	fmt.Fprintf(w, `<img src="%s" alt="%s"><br>`, foundArtist.Image, foundArtist.Name)
// 	fmt.Fprintf(w, "<h2>Members:</h2><ul>")
// 	for _, member := range foundArtist.Members {
// 		fmt.Fprintf(w, "<li>%s</li>", member)
// 	}
// 	fmt.Fprintf(w, "</ul>")
// }

func main() {
	fmt.Println("http://localhost:8080")

	http.HandleFunc("/", HomePage)
	http.HandleFunc("/artist", ArtistPage)

	http.ListenAndServe(":8080", nil)
}
