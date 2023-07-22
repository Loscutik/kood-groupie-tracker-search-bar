package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"groupietracker/api"
	"groupietracker/filters"
)

// The handler of the main page
func (app *application) homePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.NotFound(w, r)
		return
	}

	var filtersValue filtersValues
	var filtersSet filters.SetGivenFilter
	var filtersOrSet filters.SetGivenFilterOr

	// getting filters
	filtersValue = getFiltersValuesFromQuery(r.URL.Query())

	filtersSet.AddGivenFilter(filters.FilterNameContain, filtersValue.Name)

	filtersOrSet.AddGivenFilterOr(filters.FilterNumberOfMembers, filtersValue.Members)

	err := addDateCreationFilter(&filtersValue, &filtersSet)
	if err != nil {
		app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("running a filter Date Creation failures: %v ", err))
		return
	}

	err = addDateAlbumFilters(&filtersSet, "Gt", filtersValue.FirstAlbumFrom)
	if err != nil {
		app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("running a filter for First Album Date failures: %v ", err))
		return
	}

	err = addDateAlbumFilters(&filtersSet, "Lt", filtersValue.FirstAlbumTo)
	if err != nil {
		app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("running a filter for First Album Date failure: %v ", err))
		return
	}

	filtersSet.AddGivenFilter(filters.FilterLocationContain, filtersValue.Location)
	
	filtersSet.AddGivenFilter(filters.FilterNameOfMemberContain, filtersValue.NameOfMember)

	// apply all filters to set of artists
	output := &struct {
		Artists           []*api.Artist
		FiltersConstrains *filtersConstrains
		FiltersValues     *filtersValues
	}{
		Artists:           filters.ApplyFilters(app.artists, filtersSet, filtersOrSet),
		FiltersConstrains: app.filtersConstrains,
		FiltersValues:     &filtersValue,
	}

	// Assembling the page from templates
	GetTemplate("home.page.tmpl", w, output)
}

// The handler of the website data
func (app *application) concertsInfo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	app.infoLog.Printf("concertInfo - requested id: %s", id)
	i, err := strconv.Atoi(id)
	if err != nil || i < 1 || i > len(app.artists) {
		app.errLog.Printf("wrong artist's id (%s), %s", id, err)
		app.NotFound(w, r)
		return
	}
	rel, err := api.GetArtistsRelation(app.client, (app.artists)[i-1])
	if err != nil {
		app.errLog.Printf("error during getting relations %s", err)
		app.NotFound(w, r)
		return
	}
	app.infoLog.Printf("relations for id %s is: %#v\n", id, rel)
	if rel.Id != i {
		app.ServerError(w, r, "relations' Id is different from the artist's Id", fmt.Errorf("%d != %d", rel.Id, i))
		return
	}
	// Parses locations and the dates of concerts in each location
	type datesLocations []struct {
		Location string
		Dates    []string
	}
	datesLocs := make(datesLocations, len(rel.DatesLocations))
	j := 0
	for l, d := range rel.DatesLocations {
		l = strings.ReplaceAll(l, "_", " ")
		datesLocs[j].Location = strings.ReplaceAll(l, "-", ", ")
		datesLocs[j].Dates = d
		j++
	}
	output := struct {
		Id             int
		Image          string
		Name           string
		DatesLocations *datesLocations
	}{(app.artists)[i-1].Id, (app.artists)[i-1].Image, (app.artists)[i-1].Name, &datesLocs}
	GetTemplate("info.page.tmpl", w, output)
}

func getJsonForJS(w http.ResponseWriter, r *http.Request) {
	// Send GET request to the URL
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		app.ServerError(w, r, "get api failed", err)
		return
	}
	defer response.Body.Close()
	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		app.ServerError(w, r, "get api failed", err)
		return
	}

	w.Header().Add("Access-Control-Allow-Origin", "http://www.localhost:8080")
	// Print the JSON text
	jsonData := string(body) // The string of the Artists Json
	fmt.Fprint(w, jsonData)
}

func getJsonForJSloc(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		app.ServerError(w, r, "get api failed", err)
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		app.ServerError(w, r, "get api failed", err)
		return
	}

	w.Header().Add("Access-Control-Allow-Origin", "http://www.localhost:8080")
	jsonData := string(body) // The string of the Artists Json
	fmt.Fprint(w, jsonData)
}

type middlewareFunc func(http.Handler) http.Handler

func (app *application) methodChecker(methods ...string) middlewareFunc {
	return func(h http.Handler) http.Handler {
		checkMethods := func(w http.ResponseWriter, r *http.Request) {
			if methods == nil {
				h.ServeHTTP(w, r)
				return
			}
			for _, method := range methods {
				if r.Method == method {
					h.ServeHTTP(w, r)
					return
				}
			}

			app.MethodNotAllowed(w, r, methods...)
		}

		return http.HandlerFunc(checkMethods)
	}
}
