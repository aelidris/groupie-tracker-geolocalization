package Music

import (
	"net/http"
	"text/template"
)

func renderErrorPage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	errorData := ErrorData{
		StatusCode: statusCode,
		Message:    message,
	}
	tmpl, err := template.ParseFiles("views/error.html")
	if err != nil {
		http.Error(w, "Error loading error page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, errorData)
	if err != nil {
		http.Error(w, "Error rendering error page", http.StatusInternalServerError)
		return
	}
}
