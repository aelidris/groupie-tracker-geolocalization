package Music

import (
	"net/http"
	"text/template"
)

func renderErrorPage(w http.ResponseWriter, statusCode int, message string) {
	errorData := ErrorData{
		StatusCode: statusCode,
		Message:    message,
	}
	tmpl, err := template.ParseFiles("views/error.html")
	if err != nil {
		http.Error(w, "Error: loading error page", http.StatusInternalServerError)
		return
	}
	// This position of WriteHeader ensures the program only writes to the header once when the error template doesn't exist.
	w.WriteHeader(statusCode)
	err = tmpl.Execute(w, errorData)
	if err != nil {
		http.Error(w, "Error: rendering error page", http.StatusInternalServerError)
		return
	}
}
