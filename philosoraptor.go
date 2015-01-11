package main

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", handler)
	http.Handle("/", router)

	http.ListenAndServe(":8000", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}
