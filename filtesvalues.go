package main

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"groupietracker/filters"
)

type date struct {
	year, month, day string
}

type filtersValues struct {
	Name,
	CreationDateFrom,
	CreationDateTo,
	Location string
	FirstAlbumFrom,
	FirstAlbumTo date
	Members []int
	NameOfMember string
}

func checkStringDate(date string, from, to int) bool {
	dateInt, err := strconv.Atoi(date)
	return err == nil && dateInt >= from && dateInt <= to
}

func (d *date) notEmptyYear() bool {
	return d.year != ""
}

func (d *date) notEmptyMonth() bool {
	return d.month != "" && d.month != "0"
}

func (d *date) notEmptyDay() bool {
	return d.day != "" && d.day != "0"
}

func (d *date) checkYear(from, to int) bool {
	return checkStringDate(d.year, from, to)
}

func (d *date) isValidMonth() bool {
	return checkStringDate(d.month, 1, 12)
}

func (d *date) isValidDay() bool {
	return checkStringDate(d.day, 1, 31)
}

func (d *date) mergeMonthYear() string {
	return d.year + d.month
}

func (d *date) mergeDayMonthYear() string {
	return d.year + d.month + d.day
}

/*
runs a filter for the given date.
It checks if the date consists of the year, month, day, or of the year and month, or of only a year.
Then it runs the filter for the given operation. The operation must be a one of the follow strings: "Eq", "Lt", "Gt".
*/
func addDateAlbumFilters(filtersSet *filters.SetGivenFilter, operation string, currentFilter date) error {
	operationFuncs := map[string](map[string]filters.Filter){
		"Date": {
			"Eq": filters.FilterFirstAlbumDateEq,
			"Lt": filters.FilterFirstAlbumDateLt,
			"Gt": filters.FilterFirstAlbumDateGt,
		},
		"Month": {
			"Eq": filters.FilterFirstAlbumMonthEq,
			"Lt": filters.FilterFirstAlbumMonthLt,
			"Gt": filters.FilterFirstAlbumMonthGt,
		},
		"Year": {
			"Eq": filters.FilterFirstAlbumYearEq,
			"Lt": filters.FilterFirstAlbumYearLt,
			"Gt": filters.FilterFirstAlbumYearGt,
		},
	}

	if currentFilter.notEmptyYear() {
		if !currentFilter.checkYear(1900, time.Now().Year()) {
			return fmt.Errorf("wrong years for album creation %s", currentFilter.year)
		}

		if currentFilter.notEmptyMonth() {
			if !currentFilter.isValidMonth() {
				return fmt.Errorf("wrong months for album creation %s.%s", currentFilter.month, currentFilter.year)
			}

			if currentFilter.notEmptyDay() {
				if !currentFilter.isValidDay() {
					return fmt.Errorf("wrong day for album creation %s.%s.%s", currentFilter.day, currentFilter.month, currentFilter.year)
				}

				filtersSet.AddGivenFilter(operationFuncs["Date"][operation], currentFilter.mergeDayMonthYear())
			} else {
				filtersSet.AddGivenFilter(operationFuncs["Month"][operation], currentFilter.mergeMonthYear())
			}
		} else {
			filtersSet.AddGivenFilter(operationFuncs["Year"][operation], currentFilter.year)
		}
	}
	return nil
}

func addDateCreationFilter(currentFiltersValue *filtersValues, filtersSet *filters.SetGivenFilter) error {
	yearFrom := date{year: currentFiltersValue.CreationDateFrom}
	yearTo := date{year: currentFiltersValue.CreationDateTo}
	if currentFiltersValue.CreationDateFrom == "" {
		if currentFiltersValue.CreationDateTo != "" {
			if !yearTo.checkYear(1900, time.Now().Year()) {
				return fmt.Errorf("wrong years for a group creation %s", yearTo)
			}

			filtersSet.AddGivenFilter(filters.FilterYearCreatingLt, currentFiltersValue.CreationDateTo)
		}
	} else /*filtersValue.CreationDateFrom !=""*/ if currentFiltersValue.CreationDateTo == "" {
		filtersSet.AddGivenFilter(filters.FilterYearCreatingGt, currentFiltersValue.CreationDateFrom)
	} else { /*filtersValue.CreationDateFrom !="" && filtersValue.CreationDateTo !="" */

		if !yearFrom.checkYear(1900, time.Now().Year()) {
			return fmt.Errorf("wrong years for a group creation %s", yearFrom)
		}

		if !yearTo.checkYear(1900, time.Now().Year()) {
			return fmt.Errorf("wrong years for a group creation %s", yearTo)
		}

		switch {
		case currentFiltersValue.CreationDateFrom == currentFiltersValue.CreationDateTo:
			filtersSet.AddGivenFilter(filters.FilterYearCreatingEq, currentFiltersValue.CreationDateFrom)
		case currentFiltersValue.CreationDateFrom < currentFiltersValue.CreationDateTo:
			filtersSet.AddGivenFilter(filters.FilterYearCreatingGt, currentFiltersValue.CreationDateFrom)
			filtersSet.AddGivenFilter(filters.FilterYearCreatingLt, currentFiltersValue.CreationDateTo)
		}
	}
	return nil
}

func getFiltersValuesFromQuery(query url.Values) filtersValues {
	var filtersValue filtersValues
	filtersValue.Name = query.Get("Name")

	//filtersValue.Members = make([]int, len(app.filtersConstrains.Members))
	if query.Get("Members") != "" {
		membersFilters := query["Members"]
		for _, mf := range membersFilters {
			i, err := strconv.Atoi(mf)
			if err == nil {
				filtersValue.Members = append(filtersValue.Members, i)
			}
		}
	}

	filtersValue.CreationDateFrom = query.Get("CreationDateFrom")
	filtersValue.CreationDateTo = query.Get("CreationDateTo")

	filtersValue.FirstAlbumFrom.year = query.Get("FirstAlbumFromYear")
	filtersValue.FirstAlbumFrom.month = query.Get("FirstAlbumFromMonth")
	filtersValue.FirstAlbumFrom.day = query.Get("FirstAlbumFromDay")
	filtersValue.FirstAlbumTo.year = query.Get("FirstAlbumToYear")
	filtersValue.FirstAlbumTo.month = query.Get("FirstAlbumToMonth")
	filtersValue.FirstAlbumTo.day = query.Get("FirstAlbumToDay")

	filtersValue.Location = query.Get("Location")

	filtersValue.NameOfMember = query.Get("NameOfMember")

	return filtersValue
}
