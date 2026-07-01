// handlers/handlers.go
// Package handlers provides HTTP handler functions for serving the home page and artist detail pages
// in the Groupie Tracker application. It loads templates and interacts with the models package to
// retrieve artist data and render responses.
package handlers

import (
	"groupieTracker/models"
	"html/template"
	"net/http"
	"strconv"
)

type ArtistPageData struct {
	Artist    models.Artists
	Relations models.Relations
}

const (
	indexTemplatePath  = "templates/index.html"
	artistTemplatePath = "templates/artist.html"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(indexTemplatePath)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	artists := models.GetArtists()
	err = tmpl.Execute(w, artists)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing artist ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	artists := models.GetArtists()
	var selectedArtist models.Artists
	for _, artist := range artists {
		if artist.ID == id {
			selectedArtist = artist
			break
		}
	}
	if selectedArtist.ID == 0 {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	relations := models.GetArtistByID(id)

	data := ArtistPageData{
		Artist:    selectedArtist,
		Relations: relations,
	}

	artistTmpl, err := template.ParseFiles(artistTemplatePath)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	err = artistTmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

}
