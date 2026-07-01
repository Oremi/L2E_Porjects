// models/models.go
package models

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Artists struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	RelationsURL string   `json:"relations"` // URL to relations
}

type Locations struct {
	LastLocations     []string `json:"lastLocations"`
	UpcomingLocations []string `json:"upcomingLocations"`
}

type Dates struct {
	LastConcert     string `json:"lastConcert"`
	UpcomingConcert string `json:"upcomingConcert"`
}

type Relations struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// GetArtists fetches the list of artists from the external API and returns them as a slice of Artists structs.
func GetArtists() []Artists {
	var artists []Artists
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&artists); err != nil {
		log.Fatal(err)
	}

	return artists
}

// GetArtistByID fetches the details of a specific artist by ID from the external API and returns it as a Relations struct.
func GetArtistByID(id int) Relations {
	var relations Relations
	url := "https://groupietrackers.herokuapp.com/api/relation/" + strconv.Itoa(id)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(body, &relations); err != nil {
		log.Fatal(err)
	}

	return relations
}
