package Music

import (
	"html/template"
	"net/http"
)

// HomePage displays a welcome message and a list of artists
func HomePage(w http.ResponseWriter, r *http.Request) {
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

	// // Display the homepage with links to artist pages
	// fmt.Fprintf(w, "<h1>Welcome to Groupie Trackers</h1>")
	// fmt.Fprintf(w, "<h2>Artists List</h2><ul>")
	// for _, artist := range artists {
	// 	fmt.Fprintf(w, `<li><a href="/artist?id=%d">%s</a></li>`, artist.ID, artist.Name)
	// }
	// fmt.Fprintf(w, "</ul>")
}
