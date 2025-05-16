
package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	groupietracker "groupie-tracker/fetching"
)

// renderStatusPage is a helper function to render the status.html template with a message and HTTP status code
func renderStatusPage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	statusTemplate, err := template.ParseFiles("templates/status.html")
	if err != nil {
		// If template parsing fails, fallback to plain text error
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = statusTemplate.Execute(w, message)
	if err != nil {
		// If template execution fails, fallback to plain text error
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

// safeStaticHandler securely serves static files, returning a 404 page if file not found or path is invalid
func safeStaticHandler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Full path to the requested file
		path := filepath.Join(dir, r.URL.Path[len("/static/"):])

		// Check if file exists and is not a directory
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			renderStatusPage(w, http.StatusNotFound, "404 Page Not Found")
			return
		}

		// Serve the static file
		http.ServeFile(w, r, path)
	}
}

// safeImageHandler securely serves image files, returning a 404 page if file not found or path is invalid
func safeImageHandler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Full path to the requested image file
		path := filepath.Join(dir, r.URL.Path[len("/images/"):])

		// Check if file exists and is not a directory
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			renderStatusPage(w, http.StatusNotFound, "404 Page Not Found")
			return
		}

		// Serve the image file
		http.ServeFile(w, r, path)
	}
}

func main() {
	// URL to fetch artists data from external API
	url := "https://groupietrackers.herokuapp.com/api/artists"

	// Fetch artists data
	artists, err := groupietracker.FetchArtists(url)
	if err != nil {
		log.Fatalf("Failed to fetch artists: %s", err)
	}

	// Set up handlers for serving static files and images securely
	http.HandleFunc("/static/", safeStaticHandler("static"))
	http.HandleFunc("/images/", safeImageHandler("images"))

	// Root handler to display list of artists with optional search query
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Root handler called with URL path: %s", r.URL.Path)

		// Only allow root path "/"
		if r.URL.Path != "/" {
			renderStatusPage(w, http.StatusNotFound, "404 Page Not Found")
			return
		}

		// Only allow GET method
		if r.Method != http.MethodGet {
			renderStatusPage(w, http.StatusBadRequest, "400 Bad Request")
			return
		}

		// Get search query parameter
		query := r.URL.Query().Get("search")

		// Filter artists by name based on search query
		filteredArtists := groupietracker.SearchArtistsByName(artists, query)

		// Prepare data for template rendering
		data := groupietracker.PageData{
			Title:       "Groupie-Tracker",
			Artists:     filteredArtists,
			SearchQuery: query,
		}

		// Parse and execute the index template with data
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

	// Artist detail handler to display details of a specific artist
	http.HandleFunc("/artist/", func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET method
		if r.Method != http.MethodGet {
			renderStatusPage(w, http.StatusBadRequest, "400 Bad Request")
			return
		}

		// Extract artist ID from URL path
		artistID := r.URL.Path[len("/artist/"):]
		log.Printf("Artist ID requested: %s", artistID)

		// Fetch artist details by ID
		artist, err := groupietracker.FetchArtistDetails(artistID)
		if err != nil {
			log.Printf("Error fetching artist details: %v", err)
			renderStatusPage(w, http.StatusBadRequest, "400 Bad Request")
			return
		}

		// Prepare data for template rendering
		data := groupietracker.ArtistPageData{
			Title:  artist.Name,
			Artist: artist,
		}

		// Parse and execute the artist template with data
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

	// Start the HTTP server on port 8080
	log.Println("Server started on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
