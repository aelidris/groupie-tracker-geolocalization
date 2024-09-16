package Music

import (
	"html/template"
	"net/http"
)

// HomePage displays a welcome message and a list of artists
func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderErrorPage(w, http.StatusNotFound, "Page Not Found")
		return
	}
	artists, err := FetchArtists()
	if err != nil {
		http.Error(w, "Failed to load artists", http.StatusInternalServerError)
		return
	}

	temp, err := template.ParseFiles("views/index.html")
	if err != nil {
		renderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = temp.Execute(w, artists)
	if err != nil {
		renderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}
