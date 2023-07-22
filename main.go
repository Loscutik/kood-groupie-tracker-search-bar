package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"groupietracker/api"
)

// constrains for filters (output on the page)
type filtersConstrains struct {
	MinYearCreation,
	MaxYearCreation int
	Members   []int
	Locations map[string]bool
}
type application struct {
	errLog            *log.Logger
	infoLog           *log.Logger
	client            *http.Client
	apies             api.Api
	artists           []*api.Artist
	filtersConstrains *filtersConstrains
}

var app application

func main() {
	var err error
	app.createLog()
	app.createWebClient()

	// Gets data from the given API
	app.apies, err = api.GetAPI(app.client, api.API)
	if err != nil {
		log.Fatalf("fail creating api: %s", err)
	}
	app.infoLog.Printf("apies: %#v\n", app.apies)

	// Get data of the specific Artists list
	app.artists, err = api.GetArtists(app.client, app.apies)
	if err != nil {
		log.Fatalf("fail fetcing artists' data: %s", err)
	}
	for _,artist:=range app.artists{
		app.infoLog.Printf("artists: %#v\n", *artist)
	}

	app.setFilterConstrains()

	mux := app.createRouters()

	// Starting the web server
	port, err := parseArgs()
	if err != nil {
		app.errLog.Fatal(err)
	}

	fmt.Printf("Starting server at http://localhost:%s\n", *port)
	app.infoLog.Printf("Starting server at http://localhost:%s\n", *port)
	if err := http.ListenAndServe(":"+*port, mux); err != nil {
		app.errLog.Fatal(err)
	}
}

// this type doesn't allow FileServer to open directions in the static direction
type neuteredFileSystem struct {
	fs http.FileSystem
}

// implements FileSystem interface
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		if _, err := nfs.fs.Open(filepath.Join(path, "index.html")); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}

	return f, nil
}

// creates 2 logs: for information and for errors
func (app *application) createLog() {
	// Creates logs of what happenes
	errLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)                          // Creates logs of errors
	infoLogFile, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o664) // Puts log info into the specific file
	if err != nil {
		errLog.Printf("Cannot open a log file. Error is %s\nStdout will be used for the info log ", err)
		infoLogFile = os.Stdout
	}
	infoLog := log.New(infoLogFile, "INFO:  ", log.Ldate|log.Ltime|log.Lshortfile)

	app.errLog = errLog
	app.infoLog = infoLog
}

// creates a web client with 20sec timeout
func (app *application) createWebClient() {
	// Specifies the time limits for server requests
	app.client = &http.Client{
		Timeout: time.Second * 20,
	}
}

// sets constraints for filters: list of locations, max & min years of the bands creation, max number of members of a band
func (app *application) setFilterConstrains() {
	// preparing data for filters
	app.filtersConstrains = &filtersConstrains{
		MinYearCreation: 3000,
		MaxYearCreation: 0,
		Locations:       make(map[string]bool),
	}

	maxMembers := 0
	for _, a := range app.artists {
		// Get a list of locations for the filter
		for _, l := range a.Locations.Locations {
			if _, ok := app.filtersConstrains.Locations[l]; !ok {
				app.filtersConstrains.Locations[l] = false
			}
		}

		// Get max & min year cration for bands
		if a.CreationDate < app.filtersConstrains.MinYearCreation {
			app.filtersConstrains.MinYearCreation = a.CreationDate
		} else {
			if a.CreationDate > app.filtersConstrains.MaxYearCreation {
				app.filtersConstrains.MaxYearCreation = a.CreationDate
			}
		}

		// Get max number of members of the bands
		if len(a.Members) > maxMembers {
			maxMembers = len(a.Members)
		}
	}
	app.infoLog.Printf("locations' list: %#v\n", app.filtersConstrains.Locations)

	app.filtersConstrains.Members = make([]int, maxMembers)
	for i := 0; i < maxMembers; i++ {
		app.filtersConstrains.Members[i] = i + 1
	}
	app.infoLog.Printf("filters: %#v\n", app.filtersConstrains)
}

// creates a serverMux and assign handlers to it
func (app *application) createRouters() *http.ServeMux {
	mux := http.NewServeMux()
	// Handlers to run the web pages
	mux.HandleFunc("/", app.homePageHandler)
	mux.HandleFunc("/info", app.concertsInfo)
	mux.Handle("/json", app.methodChecker("POST")(http.HandlerFunc(app.getJsonForJS)))
	//mux.Handle("/jsonloc", app.methodChecker("POST")(http.HandlerFunc(getJsonForJSloc)))
	fileServer := http.FileServer(neuteredFileSystem{http.Dir(TEMPLATES_PATH + "static/")})
	// fileServer := http.FileServer( http.Dir(TEMPLATES_PATH + "static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	return mux
}

// Parses the program's arguments to obtain the server port. If no arguments found, it uses the 8080 port by default
// Usage: go run .  --port=PORT_NUMBER
func parseArgs() (*string, error) {
	port := flag.String("port", "8080", "server port")
	flag.Parse()
	if flag.NArg() > 0 {
		return nil, fmt.Errorf("wrong arguments\nUsage: go run .  --port=PORT_NUMBER")
	}
	_, err := strconv.ParseUint(*port, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("error: port must be a 16-bit unsigned number ")
	}
	return port, nil
}
