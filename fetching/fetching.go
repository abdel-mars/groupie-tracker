package groupietracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// Artist represents the basic information about a music artist.
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

// ArtistDetails extends Artist with additional detailed information.
type ArtistDetails struct {
	Artist
	LocationsList    []string            `json:"locations"`
	ConcertDatesList []string            `json:"concertDates"`
	RelationsMap     map[string][]string `json:"relations"`
}

// PageData holds data for rendering the main artists list page.
type PageData struct {
	Title       string
	Artists     []Artist
	SearchQuery string
}

// ArtistPageData holds data for rendering an individual artist's detail page.
type ArtistPageData struct {
	Title  string
	Artist ArtistDetails
}

// FetchArtists fetches the list of artists from the given URL.
func FetchArtists(url string) ([]Artist, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET request failed: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var artists []Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return artists, nil
}

// FetchArtistDetails fetches detailed information about an artist by ID using goroutines for concurrency.
func FetchArtistDetails(artistID string) (ArtistDetails, error) {
	var wg sync.WaitGroup
	var artist Artist
	var locations []string
	var dates []string
	var relations map[string][]string
	var err1, err2, err3, err4 error

	wg.Add(4)

	go func() {
		defer wg.Done()
		artist, err1 = FetchArtist(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%s", artistID))
	}()

	go func() {
		defer wg.Done()
		locations, err2 = FetchLocations(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%s", artistID))
	}()

	go func() {
		defer wg.Done()
		dates, err3 = FetchDates(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/dates/%s", artistID))
	}()

	go func() {
		defer wg.Done()
		relations, err4 = FetchRelations(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%s", artistID))
	}()

	wg.Wait()

	if err1 != nil {
		return ArtistDetails{}, err1
	}
	if err2 != nil {
		return ArtistDetails{}, err2
	}
	if err3 != nil {
		return ArtistDetails{}, err3
	}
	if err4 != nil {
		return ArtistDetails{}, err4
	}

	details := ArtistDetails{
		Artist:           artist,
		LocationsList:    locations,
		ConcertDatesList: dates,
		RelationsMap:     relations,
	}

	return details, nil
}

// FetchArtist fetches basic artist information from the given URL.
func FetchArtist(url string) (Artist, error) {
	response, err := http.Get(url)
	if err != nil {
		return Artist{}, fmt.Errorf("HTTP GET request failed: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Artist{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var artist Artist
	err = json.Unmarshal(body, &artist)
	if err != nil {
		return Artist{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return artist, nil
}

// FetchLocations fetches a list of locations from the given URL.
func FetchLocations(url string) ([]string, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET request failed: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var locations struct {
		Locations []string `json:"locations"`
	}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return locations.Locations, nil
}

// FetchDates fetches a list of concert dates from the given URL.
func FetchDates(url string) ([]string, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET request failed: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var dates struct {
		Dates []string `json:"dates"`
	}
	err = json.Unmarshal(body, &dates)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return dates.Dates, nil
}

// FetchRelations fetches a map of relations from the given URL.
func FetchRelations(url string) (map[string][]string, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET request failed: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var relations struct {
		DatesLocations map[string][]string `json:"datesLocations"`
	}
	err = json.Unmarshal(body, &relations)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return relations.DatesLocations, nil
}

// SearchArtistsByName filters artists by name using case-insensitive substring matching.
func SearchArtistsByName(artists []Artist, query string) []Artist {
	if query == "" {
		return artists
	}
	var filtered []Artist
	for _, artist := range artists {
		if containsIgnoreCase(artist.Name, query) {
			filtered = append(filtered, artist)
		}
	}
	return filtered
}

// containsIgnoreCase checks if substr is a case-insensitive substring of str.
func containsIgnoreCase(str, substr string) bool {
	strLower := strings.ToLower(str)
	substrLower := strings.ToLower(substr)
	return strings.Contains(strLower, substrLower)
}
