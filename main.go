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
	http.HandleFunc("/public/", Music.SetupStaticFilesHandlers)
	http.HandleFunc("/", Music.HomePage)
	http.HandleFunc("/artist", Music.ArtistPage)
	err := http.ListenAndServe(":3040", nil)
	if err != nil {
		log.Fatal(err)
	}
}
