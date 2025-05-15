package main

import (
	"html/template"
	"log"
	"net/http"
	groupietracker "groupie-tracker/fetching"
)

// main initializes the server, fetches artist data, and sets up HTTP handlers.
func main() {
	// URL to fetch the list of artists
	url := "https://groupietrackers.herokuapp.com/api/artists"

	// Fetch artists data from the API
	artists, err := groupietracker.FetchArtists(url)
	if err != nil {
		log.Fatalf("Failed to fetch artists: %s", err)
	}

	// Serve static files from the ./static directory
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve images from the ./images directory
	imgFs := http.FileServer(http.Dir("./images"))
	http.Handle("/images/", http.StripPrefix("/images/", imgFs))

	// Handle the root path "/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Root handler called with URL path: %s", r.URL.Path)

		// Parse the status template for error handling
		statusTemplate, err := template.ParseFiles("templates/status.html")
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Only allow root path "/"
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			statusTemplate.Execute(w, "404 Page Not Found")
			return
		}

		// Only allow GET method
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			statusTemplate.Execute(w, "400 Bad Request")
			return
		}

		// Get search query parameter
		query := r.URL.Query().Get("search")

		// Filter artists by search query
		filteredArtists := groupietracker.SearchArtistsByName(artists, query)

		// Prepare data for the template
		data := groupietracker.PageData{
			Title:       "Artists List",
			Artists:     filteredArtists,
			SearchQuery: query,
		}

		// Parse the index template
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			statusTemplate.Execute(w, "500 Internal Server Error")
			return
		}

		// Execute the template with data
		err = tmpl.Execute(w, data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			statusTemplate.Execute(w, "500 Internal Server Error")
			return
		}
	})

	// Handle artist detail pages "/artist/{id}"
	http.HandleFunc("/artist/", func(w http.ResponseWriter, r *http.Request) {
		// Parse the status template for error handling
		statusTemplate, err := template.ParseFiles("templates/status.html")
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Only allow GET method
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			statusTemplate.Execute(w, "400 Bad Request")
			return
		}

		// Extract artist ID from URL path
		artistID := r.URL.Path[len("/artist/"):]
		log.Printf("Artist ID requested: %s", artistID)

		// Fetch artist details by ID
		artist, err := groupietracker.FetchArtistDetails(artistID)
		if err != nil {
			log.Printf("Error fetching artist details: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			statusTemplate.Execute(w, "400 Bad Request")
			return
		}

		// Prepare data for the template
		data := groupietracker.ArtistPageData{
			Title:  artist.Name,
			Artist: artist,
		}

		// Parse the artist template
		tmpl, err := template.ParseFiles("templates/artist.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			statusTemplate.Execute(w, "500 Internal Server Error")
			return
		}

		// Execute the template with data
		err = tmpl.Execute(w, data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			statusTemplate.Execute(w, "500 Internal Server Error")
			return
		}
	})

	log.Println("Server started on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
