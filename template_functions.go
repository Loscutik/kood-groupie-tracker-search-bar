package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

const (
	TEMPLATES_PATH = "./templates/"
)

func GetTemplate(a string, w http.ResponseWriter, s any) {
	site := []string{
		TEMPLATES_PATH + "base.layout.tmpl",
		TEMPLATES_PATH + a,
	}

	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"nowYear": time.Now().Year,
	}
	tm, err := template.New("base.layout.tmpl").Funcs(funcMap).ParseFiles(site...)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}
	err = tm.Execute(w, s)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}
}

// Opens a beautiful HTML 404 web page instead of the status 404 "Page not found"
func (app *application) NotFound(w http.ResponseWriter, r *http.Request) {
	app.errLog.Printf("wrong path: %s", r.URL.Path)
	w.WriteHeader(http.StatusNotFound)                             // Sets status code at 404
	tm, _ := template.ParseFiles(TEMPLATES_PATH + "error404.html") // Opens the HTML web page
	err := tm.Execute(w, nil)
	if err != nil {
		app.errLog.Println(err)
		http.NotFound(w, r)
	}
}

func (app *application) ServerError(w http.ResponseWriter, r *http.Request, message string, err error) {
	app.errLog.Output(2, fmt.Sprintf("fail handling the page %s: %s: %s\n%v", r.URL.Path, message, err, debug.Stack()))
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func (app *application) ClientError(w http.ResponseWriter, r *http.Request, errStatus int, logTexterr string) {
	app.errLog.Output(2, logTexterr)
	http.Error(w, "ERROR: "+http.StatusText(errStatus), errStatus)
}

func (app *application) MethodNotAllowed(w http.ResponseWriter, r *http.Request, allowedMethods ...string) {
	if allowedMethods == nil {
		panic("no methods is given to func MethodNotAllowed")
	}
	allowdeString := allowedMethods[0]
	for i := 1; i < len(allowedMethods); i++ {
		allowdeString += ", " + allowedMethods[i]
	}

	w.Header().Set("Allow", allowdeString)
	app.ClientError(w, r, http.StatusMethodNotAllowed, fmt.Sprintf("using the method %s to go to a page %s", r.Method, r.URL))
}
