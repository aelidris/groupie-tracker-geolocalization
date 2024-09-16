package Music

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)



// GetApi fetches data from a given API URL and populates the Urls struct
func GetApi(api string) {
	response, err := http.Get(api) // Send an HTTP GET request to the URL=api to fetch the data
	if err != nil {
		log.Fatal(err) // Handle errors
	}
	// Decode the data from JSON
	json.NewDecoder(response.Body).Decode(&Urls)
	defer response.Body.Close() // Close resources
}

// GetArtists handles requests for both the homepage and artist-specific pages
func GetArtists(w http.ResponseWriter, r *http.Request) {
	// Serve error template in case of errors
	tmp, err := template.ParseFiles("views/error.html")
	if err != nil {
		http.Error(w, "Error page not found: Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch artist data if not already populated
	if len(Artist) == 0 {
		artistUrl := Urls.ArtistsUrl
		artistRes, err := http.Get(artistUrl)
		if err != nil {
			renderErrorPage(w, http.StatusInternalServerError, "Failed to fetch artist data")
			return
		}
		defer artistRes.Body.Close()

		err = json.NewDecoder(artistRes.Body).Decode(&Artist)
		if err != nil {
			renderErrorPage(w, http.StatusInternalServerError, "Failed to decode artist data")
			return
		}
	}

	// Check if the request is for the homepage or an artist-specific page
	if r.URL.Path == "/" {
		// Serve the homepage
		tmp, err = template.ParseFiles("views/index.html")
		if err != nil {
			renderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		err = tmp.Execute(w, Artist) // Pass artist data to the homepage
		if err != nil {
			renderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	} else if strings.Contains(r.URL.Path, "/artist/") {
		// Serve artist-specific page
		StrId := r.URL.Path[len("/artist/"):]
		Id, err := strconv.Atoi(StrId)
		if err != nil || Id < 1 || Id > len(Artist) || strings.HasPrefix(StrId, "0") {
			renderErrorPage(w, http.StatusNotFound, "Invalid artist ID")
			return
		}

		FetchArtistData(Id, w) // Fetch and render the artist's data
	} else {
		// Handle other routes as 404
		renderErrorPage(w, http.StatusNotFound, "Page Not Found")
		return
	}
}

// FetchArtistData fetches data for a specific artist and renders the artist.html template
func FetchArtistData(Id int, w http.ResponseWriter) {
	// Prepare the map data to pass to the template
	data := MapData{
		APIKey: "84a9a08245a141e299e5a1fd45b3dbd8", // Your Geoapify API key
		Loca: [][2]float64{
			{48.8588443, 2.2943506},   // Eiffel Tower, Paris
			{40.6892494, -74.0445004}, // Statue of Liberty, New York
			{35.658581, 139.745438},   // Tokyo Tower, Tokyo
			{51.5007292, -0.1268141},  // Big Ben, London
		},
	}

	// Adjust for zero-based indexing
	Id -= 1

	// Helper function to fetch and decode JSON data
	fetchAndDecode := func(url string, target interface{}) error {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return json.NewDecoder(resp.Body).Decode(target)
	}

	// Fetch location, date, and relations data for the artist
	err := fetchAndDecode(Artist[Id].LocationsURL, &Artist[Id].Location)
	if err != nil {
		renderErrorPage(w, http.StatusInternalServerError, "Failed to fetch location data")
		return
	}

	err = fetchAndDecode(Artist[Id].ConcertDatesURL, &Artist[Id].Date)
	if err != nil {
		renderErrorPage(w, http.StatusInternalServerError, "Failed to fetch concert dates")
		return
	}

	err = fetchAndDecode(Artist[Id].RelationsURL, &Artist[Id].Relation)
	if err != nil {
		renderErrorPage(w, http.StatusInternalServerError, "Failed to fetch relations data")
		return
	}

	// Parse and execute the artist-specific template
	tmp, err := template.ParseFiles("views/artist.html")
	if err != nil {
		renderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Execute the template, passing the specific artist's data
	err = tmp.Execute(w, Artist[Id])
	if err != nil {
		renderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

