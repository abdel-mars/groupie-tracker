# Groupie Tracker

![Project Logo](images/plogo.png)
visit link : https://groupie-tracker-ykss.onrender.com

Groupie Tracker is a web application built with Go that allows users to explore information about music artists. It fetches artist data from an external API and provides a user-friendly interface to browse, search, and view detailed information about artists, including their members, locations, concert dates, and relations.

## Features

- Display a list of music artists with images and basic info.
- Search artists by name with case-insensitive filtering.
- View detailed artist information including:
  - Members of the band
  - Locations where the artist has performed
  - Concert dates
  - Relations with other artists
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
│   ├── style-dark.css
│   └── style-404.css
├── images/                 # Image assets used in the UI
│   └── plogo.png
└── templates/              # HTML templates for rendering pages
    ├── index.html
    ├── artist.html
    └── status.html
```

## API Endpoints Used

- `https://groupietrackers.herokuapp.com/api/artists` - Fetches the list of all artists.
- `https://groupietrackers.herokuapp.com/api/artists/{id}` - Fetches detailed information for a specific artist by ID.
- `https://groupietrackers.herokuapp.com/api/locations/{id}` - Fetches locations for a specific artist.
- `https://groupietrackers.herokuapp.com/api/dates/{id}` - Fetches concert dates for a specific artist.
- `https://groupietrackers.herokuapp.com/api/relation/{id}` - Fetches relations for a specific artist.


# How the Groupie Tracker Project Works

## Project Structure and Overview

I structured this project as a web application written in Go that fetches and displays information about music artists from an external API called Groupie Trackers. The project consists of:

- `main.go`: The main entry point of the web server. It handles HTTP requests, serves static files, renders HTML templates, and coordinates data fetching.
- `fetching/fetching.go`: A package I wrote to handle all API fetching and data processing logic. It fetches artist data, detailed artist info, and provides search functionality.
- `templates/`: Contains HTML templates (`index.html`, `artist.html`, `status.html`) used to render the web pages dynamically.
- `static/` and `images/`: Contain CSS stylesheets and image assets served securely by the server.

## How the Project Works Technically

### main.go

This is the core of the web server. Here is how it works:

- **Data Fetching:**  
  At startup, I fetch the list of artists from the external API using the `FetchArtists` function from the `fetching` package. This gives me the initial data to display.

- **HTTP Handlers:**  
  I set up several HTTP handlers:
  - `/static/` and `/images/` serve static CSS and image files securely, checking if files exist before serving.
  - `/` (root path) serves the main page showing the list of artists. It supports a search query parameter to filter artists by name.
  - `/artist/{id}` serves a detailed page for a specific artist by fetching detailed info concurrently from multiple API endpoints.

- **Template Rendering:**  
  I use Go's `html/template` package to parse and render HTML templates dynamically with the data fetched. This allows me to separate the HTML structure from the Go code cleanly.

- **Error Handling:**  
  I handle errors gracefully by rendering a status page with appropriate HTTP status codes and messages.

- **Server Startup:**  
  The server listens on port 8080 and logs requests and errors.

### fetching/fetching.go

This package handles all the API interactions and data processing:

- **Data Structures:**  
  I defined `Artist` and `ArtistDetails` structs to represent the artist data and detailed info.

- **Fetching Functions:**  
  I wrote functions to fetch:
  - The list of artists (`FetchArtists`)
  - Basic artist info (`FetchArtist`)
  - Locations (`FetchLocations`)
  - Concert dates (`FetchDates`)
  - Relations (`FetchRelations`)

- **Concurrency:**  
  For fetching detailed artist info, I use Go's goroutines and `sync.WaitGroup` to fetch multiple pieces of data concurrently, improving performance.

- **Search Functionality:**  
  I implemented a case-insensitive substring search function to filter artists by name.

- **Error Handling:**  
  All fetching functions handle HTTP errors and JSON parsing errors properly.

## Libraries and Code I Used or Built

- I used Go's standard library extensively:
  - `net/http` for HTTP server and client requests.
  - `html/template` for HTML templating.
  - `encoding/json` for JSON parsing.
  - `sync` for concurrency control.
  - `log` for logging.

- I wrote all the application-specific code myself, including:
  - The HTTP handlers and server setup in `main.go`.
  - The fetching and data processing logic in `fetching/fetching.go`.
  - The search functionality.
  - The secure static file serving logic.

## How HTTP and API Fetching Works (Client-Server Interaction)

- When a user visits the root URL `/`, the server handles the request by:
  - Checking for a search query parameter.
  - Filtering the list of artists accordingly.
  - Rendering the `index.html` template with the filtered data.
  - Sending the rendered HTML back to the user's browser.

- When a user clicks on an artist, the browser requests `/artist/{id}`:
  - The server fetches detailed artist info concurrently from multiple API endpoints.
  - It renders the `artist.html` template with this detailed data.
  - Sends the rendered HTML to the browser.

- Static assets like CSS and images are served securely by checking file existence before serving.

- The communication between the server and the external API is done via HTTP GET requests, fetching JSON data which is parsed into Go structs.

## Event Creation and Visualization

- The "events" in this project are the concert dates and locations associated with each artist.
- These are fetched from the API and included in the detailed artist data.
- In the `artist.html` template, I display these events visually as lists or maps (depending on the template design).
- This allows users to see where and when artists have concerts.

## Problems I Faced While Building the Project

- Handling concurrent API requests for detailed artist info required careful synchronization using goroutines and `sync.WaitGroup`.
- Ensuring secure serving of static files to prevent directory traversal attacks.
- Managing error handling gracefully to provide meaningful feedback to users.
- Parsing complex JSON structures like relations and concert dates into usable Go data structures.
- Designing templates to display dynamic data cleanly and responsively.

## Technologies Used

- Programming Language: Go (Golang)
- Web Server: net/http package
- HTML Templating: html/template package
- Concurrency: goroutines and sync package
- External API: Groupie Trackers API (https://groupietrackers.herokuapp.com/api)
- Frontend: HTML, CSS (served as static files)

---

I hope this detailed explanation helps you understand how I structured and built the Groupie Tracker project, how it works technically, and the challenges I faced. If you have any questions or want me to explain any part further, please let me know.
