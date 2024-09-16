package Music

type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Relation     Relation
	Location     Location
	Date         Date
}
type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type LatLng struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type ArtistPageData struct {
	Artist       Artist
	Location     Location
	LocationData string // String of lat/lng for passing to JavaScript
	Date Date
	Relation Relation
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type ErrorData struct {
	StatusCode int
	Message    string
}
