package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

var htmlTemplates *template.Template

func main() {
	router := mux.NewRouter()
	htmlTemplates = template.Must(template.ParseGlob("templates/*"))

	router.HandleFunc("/", homePage)
	router.HandleFunc("/generate", handleForm).Methods("POST")
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

func handleForm(w http.ResponseWriter, r *http.Request) {
	upperText := r.FormValue("upper_text")
	lowerText := r.FormValue("lower_text")
	fmt.Println(upperText)
	fmt.Println(lowerText)
	http.Redirect(w, r, "/", 301)
}
