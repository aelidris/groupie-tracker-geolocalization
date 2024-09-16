package main

import (
	"fmt"
	"log"
	"net/http"

	Music "groupietracker/hundlers"
)

const (
	api = "https://groupietrackers.herokuapp.com/api"
)

func main() {
	fmt.Println("http://localhost:3040")
	Music.GetApi(api)
	http.HandleFunc("/public/", Music.SetupStaticFilesHandlers)
	http.HandleFunc("/", Music.GetArtists)
	err := http.ListenAndServe(":3040", nil)
	if err != nil {
		log.Fatal(err)
	}
}
