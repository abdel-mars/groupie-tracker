
// Package main implements the web server for the Groupie Tracker application.
// It serves artist data fetched from an external API, renders HTML templates,
// and handles static assets securely.
package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	groupietracker "groupie-tracker/fetching"
)

// renderStatusPage renders the status.html template with a given message and HTTP status code.
// It is used to display error or status messages to the user.
func renderStatusPage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	statusTemplate, err := template.ParseFiles("templates/status.html")
	if err != nil {
		// If template parsing fails, fallback to plain text error response.
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = statusTemplate.Execute(w, message)
	if err != nil {
		// If template execution fails, fallback to plain text error response.
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

// safeStaticHandler returns an HTTP handler function that securely serves static files
// from the specified directory. If the requested file does not exist or is a directory,
// it responds with a 404 status page.
func safeStaticHandler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Construct the full file path by joining the directory and the requested URL path.
		path := filepath.Join(dir, r.URL.Path[len("/static/"):])

		// Check if the file exists and is not a directory.
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			renderStatusPage(w, http.StatusNotFound, "404 Page Not Found")
			return
		}

		// Serve the static file.
		http.ServeFile(w, r, path)
	}
}

// safeImageHandler returns an HTTP handler function that securely serves image files
// from the specified directory. If the requested image does not exist or is a directory,
// it responds with a 404 status page.
func safeImageHandler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Construct the full image file path by joining the directory and the requested URL path.
		path := filepath.Join(dir, r.URL.Path[len("/images/"):])

		// Check if the image file exists and is not a directory.
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			renderStatusPage(w, http.StatusNotFound, "404 Page Not Found")
			return
		}

		// Serve the image file.
		http.ServeFile(w, r, path)
	}
}

func main() {
	// URL of the external API to fetch artists data.
	url := "https://groupietrackers.herokuapp.com/api/artists"

	// Fetch the list of artists from the API.
	artists, err := groupietracker.FetchArtists(url)
	if err != nil {
		log.Fatalf("Failed to fetch artists: %s", err)
	}

	// Set up HTTP handlers for serving static files and images securely.
	http.HandleFunc("/static/", safeStaticHandler("static"))
	http.HandleFunc("/images/", safeImageHandler("images"))

	// Root handler to display the list of artists with optional search functionality.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Root handler called with URL path: %s", r.URL.Path)

		// Only allow the root path "/".
		if r.URL.Path != "/" {
			renderStatusPage(w, http.StatusNotFound, "404 Page Not Found")
			return
		}

		// Only allow GET requests.
		if r.Method != http.MethodGet {
			renderStatusPage(w, http.StatusBadRequest, "400 Bad Request")
			return
		}

		// Get the search query parameter from the URL.
		query := r.URL.Query().Get("search")

		// Filter artists by name based on the search query.
		filteredArtists := groupietracker.SearchArtistsByName(artists, query)

		// Prepare data for rendering the index template.
		data := groupietracker.PageData{
			Title:       "Groupie-Tracker",
			Artists:     filteredArtists,
			SearchQuery: query,
		}

		// Parse and execute the index.html template with the data.
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			renderStatusPage(w, http.StatusInternalServerError, "500 Internal Server Error")
			return
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			renderStatusPage(w, http.StatusInternalServerError, "500 Internal Server Error")
			return
		}
	})

	// Handler to display detailed information about a specific artist.
	http.HandleFunc("/artist/", func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET requests.
		if r.Method != http.MethodGet {
			renderStatusPage(w, http.StatusBadRequest, "400 Bad Request")
			return
		}

		// Extract the artist ID from the URL path.
		artistID := r.URL.Path[len("/artist/"):]
		log.Printf("Artist ID requested: %s", artistID)

		// Fetch detailed artist information by ID.
		artist, err := groupietracker.FetchArtistDetails(artistID)
		if err != nil {
			log.Printf("Error fetching artist details: %v", err)
			renderStatusPage(w, http.StatusBadRequest, "400 Bad Request")
			return
		}

		// Prepare data for rendering the artist template.
		data := groupietracker.ArtistPageData{
			Title:  artist.Name,
			Artist: artist,
		}

		// Parse and execute the artist.html template with the data.
		tmpl, err := template.ParseFiles("templates/artist.html")
		if err != nil {
			renderStatusPage(w, http.StatusInternalServerError, "500 Internal Server Error")
			return
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			renderStatusPage(w, http.StatusInternalServerError, "500 Internal Server Error")
			return
		}
	})

	// Start the HTTP server on port 8080.
	log.Println("Server started on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
