package Music

import (
	"net/http"
	"os"
)

func SetupStaticFilesHandlers(w http.ResponseWriter, r *http.Request) {
	// Construct the file path
	filePath := "." + r.URL.Path
	// Check if the file exists and is not a directory
	fileInfo, err := os.Stat(filePath)
	if err != nil || os.IsNotExist(err) || fileInfo.IsDir() {
		renderErrorPage(w, http.StatusNotFound, "Page Not Found")
		return
	}
	// Serve the file
	http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))).ServeHTTP(w, r)
}
