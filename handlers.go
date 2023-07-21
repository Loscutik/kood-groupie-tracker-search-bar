package main

import (
	"fmt"
	"io/ioutil"
	"log"
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
	// get filters from the http request

	var filtersValue filtersValues
	filtersValue.Members = make([]int, len(app.filtersConstrains.Members))
	var filtersSet []filters.SetGivenFilter
	var filtersOrSet []filters.SetGivenFilterOr

	// getting filters
	// filter by the name of groupe
	filtersValue.Name = r.URL.Query().Get("Name")
	if filtersValue.Name != "" {
		filters.AddGivenFilter(&filtersSet, filters.FilterNameContain, filtersValue.Name)
	}

	// filter by the number of members
	if r.URL.Query().Get("Members") != "" {
		membersFilters := r.URL.Query()["Members"]
		var valuesMembersFilter []int
		for _, mf := range membersFilters {
			i, err := strconv.Atoi(mf)
			if err == nil {
				valuesMembersFilter = append(valuesMembersFilter, i)
			}
		}
		for _, v := range valuesMembersFilter {
			filtersValue.Members[v-1] = 1
		}

		filters.AddGivenFilterOr(&filtersOrSet, filters.FilterNumberOfMembers, valuesMembersFilter)
	}

	// filter by the date of creation
	filtersValue.CreationDateFrom = r.URL.Query().Get("CreationDateFrom")
	filtersValue.CreationDateTo = r.URL.Query().Get("CreationDateTo")
	if filtersValue.CreationDateFrom == "" {
		if filtersValue.CreationDateTo != "" {
			filters.AddGivenFilter(&filtersSet, filters.FilterYearCreatingLt, filtersValue.CreationDateTo)
		}
	} else /*filtersValue.CreationDateFrom !=""*/ if filtersValue.CreationDateTo == "" {
		filters.AddGivenFilter(&filtersSet, filters.FilterYearCreatingGt, filtersValue.CreationDateFrom)
	} else { /*filtersValue.CreationDateFrom !="" && filtersValue.CreationDateTo !="" */
		switch {
		case filtersValue.CreationDateFrom == filtersValue.CreationDateTo:
			filters.AddGivenFilter(&filtersSet, filters.FilterYearCreatingEq, filtersValue.CreationDateFrom)
		case filtersValue.CreationDateFrom < filtersValue.CreationDateTo:
			filters.AddGivenFilter(&filtersSet, filters.FilterYearCreatingGt, filtersValue.CreationDateFrom)
			filters.AddGivenFilter(&filtersSet, filters.FilterYearCreatingLt, filtersValue.CreationDateTo)
		}
	}

	// filter by the first album date
	filtersValue.FirstAlbumFrom.year = r.URL.Query().Get("FirstAlbumFromYear")
	filtersValue.FirstAlbumFrom.month = r.URL.Query().Get("FirstAlbumFromMonth")
	filtersValue.FirstAlbumFrom.day = r.URL.Query().Get("FirstAlbumFromDay")
	filtersValue.FirstAlbumTo.year = r.URL.Query().Get("FirstAlbumToYear")
	filtersValue.FirstAlbumTo.month = r.URL.Query().Get("FirstAlbumToMonth")
	filtersValue.FirstAlbumTo.day = r.URL.Query().Get("FirstAlbumToDay")

	err := runDateFilter(filtersValue.FirstAlbumFrom, "Gt", &filtersSet)
	if err != nil {
		app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("running a filter for First Album Date failure: %v ", err))
		return
	}

	err = runDateFilter(filtersValue.FirstAlbumTo, "Lt", &filtersSet)
	if err != nil {
		app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("running a filter for First Album Date failure: %v ", err))
		return
	}

	// filter by the location
	filtersValue.Location = r.URL.Query().Get("Location")
	if filtersValue.Location != "" {
		filters.AddGivenFilter(&filtersSet, filters.FilterLocationContain, filtersValue.Location)
	}

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
		app.ServerError(w, r, "relations' Id is different from the artist's Id", err)
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

func (app *application) getJsonForJS(w http.ResponseWriter, r *http.Request) {
	// if r.Method != "POST" {
	// 	app.MethodNotAllowed(w, r, "POST")
	// 	return
	// }
	// Send GET request to the URL
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Add("Access-Control-Allow-Origin","http://www.localhost:8080")
	// Print the JSON text
	jsonData := string(body) // The string of the Artists Json
	fmt.Fprint(w, jsonData)
}

func (app *application) getJsonForJSloc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		app.MethodNotAllowed(w, r, "POST")
		return
	}
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Add("Access-Control-Allow-Origin","http://www.localhost:8080")
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
