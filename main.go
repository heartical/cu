package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.HandleFunc("/", notErrorHandler).Methods("GET")
	r.HandleFunc("/game/{room_id}", playHandler).Methods("GET")
	r.HandleFunc("/rooms", rooms).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(notErrorHandler)

	log.Println("Starting server on :8010")
	err := http.ListenAndServe(":8010", r)
	if err != nil {
		panic(err)
	}
}

func notErrorHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/game/game-not-found", http.StatusFound)
}

func rooms(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test-room"))
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["room_id"]

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Printf("error parsing template: %s", err)
		return
	}

	err = tmpl.Execute(w, roomID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Printf("error executing template: %s", err)
		return
	}
}
