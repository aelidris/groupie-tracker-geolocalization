package Music

import (
	"encoding/json"
	"net/http"
)

// FetchArtists fetches all artist data from the API
func FetchArtists() ([]Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []Artist
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}
