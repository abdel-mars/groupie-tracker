# Groupie Tracker

Groupie Tracker is a web application built with Go that allows users to explore information about music artists. It fetches artist data from an external API and provides a user-friendly interface to browse, search, and view detailed information about artists, including their members, locations, concert dates, and relations.

## Features

- [Display a list of music artists](#) with images and basic info.
- [Search artists by name](#) with case-insensitive filtering.
- [View detailed artist information](#) including:
  - [Members of the band](#)
  - [Locations where the artist has performed](#)
  - [Concert dates](#)
  - [Relations with other artists](#)
- Responsive and clean user interface using HTML templates and CSS.
- Error handling for invalid routes and requests.

## Technologies Used

- Go programming language for backend server and API integration.
- HTML templates for dynamic page rendering.
- CSS for styling (dark theme).
- External API: [Groupie Trackers API](https://groupietrackers.herokuapp.com/api/artists) for artist data.

## How to Run

1. Ensure you have [Go](https://golang.org/dl/) installed (version 1.18 or higher recommended).
2. Clone the repository and navigate to the project directory.
3. Run the application using the command:

   ```bash
   go run main.go
   ```

4. Open your web browser and visit [http://localhost:8080](http://localhost:8080) to access the app.

## Project Structure

```
.
├── main.go                 # Main application entry point
├── go.mod                  # Go module file
├── fetching/               # Package for fetching and processing artist data
│   └── fetching.go
├── static/                 # Static assets (CSS files)
│   ├── style.css
│   ├── style-dark.css
│   └── style-404.css
├── images/                 # Image assets used in the UI
│   ├── mylogo.png
│   └── plogo.png
└── templates/              # HTML templates for rendering pages
    ├── index.html
    ├── artist.html
    ├── search.html
    ├── index-err.html
    └── status.html
```

## Author

eabderrahma

---
