package main

import (
	"groupieTracker/handlers"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Test HomeHandler
func TestHomeHandler_Success(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handlers.HomeHandler(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(string(body), "<html>") {
		t.Errorf("Expected HTML content")
	}
}

// Test ArtistHandler with missing ID
func TestArtistHandler_MissingID(t *testing.T) {
	req := httptest.NewRequest("GET", "/artist", nil)
	w := httptest.NewRecorder()

	handlers.ArtistHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Missing artist ID") {
		t.Errorf("Expected 'Missing artist ID' in response")
	}
}

// Test ArtistHandler with invalid ID
func TestArtistHandler_InvalidID(t *testing.T) {
	req := httptest.NewRequest("GET", "/artist?id=abc", nil)
	w := httptest.NewRecorder()

	handlers.ArtistHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Invalid artist ID") {
		t.Errorf("Expected 'Invalid artist ID' in response")
	}
}

// Test ArtistHandler with non-existent ID
func TestArtistHandler_NonExistentID(t *testing.T) {
	req := httptest.NewRequest("GET", "/artist?id=99999", nil)
	w := httptest.NewRecorder()

	handlers.ArtistHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Artist not found") {
		t.Errorf("Expected 'Artist not found' in response")
	}
}

// Test ArtistHandler with valid ID (will depend on API)
func TestArtistHandler_ValidID(t *testing.T) {
	req := httptest.NewRequest("GET", "/artist?id=1", nil)
	w := httptest.NewRecorder()

	handlers.ArtistHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Check if template rendered
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "<html>") {
		t.Errorf("Expected HTML content")
	}
}

// Test static file serving
func TestStaticFileServing(t *testing.T) {
	// Create a test server with the mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/artist", handlers.ArtistHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/static/main.css")
	if err != nil {
		t.Fatalf("Failed to get static file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for static file, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "body {") {
		t.Errorf("Expected CSS content")
	}
}

// Test template parsing errors (simulate by removing template)
func TestTemplateError(t *testing.T) {
	// Backup original template
	originalContent, _ := os.ReadFile("templates/artist.html")
	defer os.WriteFile("templates/artist.html", originalContent, 0644)

	// Remove template file to simulate error
	os.Remove("templates/artist.html")

	req := httptest.NewRequest("GET", "/artist?id=1", nil)
	w := httptest.NewRecorder()

	handlers.ArtistHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500 for template error, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Error loading template") {
		t.Errorf("Expected 'Error loading template' in response")
	}
}
