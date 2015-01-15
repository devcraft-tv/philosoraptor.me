package main

import (
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

var htmlTemplates *template.Template

func main() {
	router := mux.NewRouter()
	htmlTemplates = template.Must(template.ParseGlob("templates/*"))

	router.HandleFunc("/", homePage)
	router.PathPrefix("/assets/").Handler(staticHandler())
	http.Handle("/", router)

	http.ListenAndServe(":8001", nil)
}

func staticHandler() http.Handler {
	return http.FileServer(http.Dir("static/"))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	err := htmlTemplates.ExecuteTemplate(w, "index", nil)
	if err != nil {
		panic(err)
	}
}
